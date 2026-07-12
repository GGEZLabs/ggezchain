#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------------------------------------------------
# estimate-upgrade-height.sh
#
# Queries a running node to estimate a safe upgrade height by:
#   1. Fetching the current block height
#   2. Calculating average block time from recent blocks
#   3. Fetching governance voting period
#   4. Estimating the block height at which voting ends
#
# Usage:
#   ./estimate-upgrade-height.sh --rpc <url> --rest <url> [--sample <n>] [--buffer <n>]
#
# Examples:
#   ./estimate-upgrade-height.sh \
#     --rpc http://localhost:26657 \
#     --rest http://localhost:1317
#
#   ./estimate-upgrade-height.sh \
#     --rpc http://localhost:26657 \
#     --rest http://localhost:1317 \
#     --sample 200 \
#     --buffer 500
# ---------------------------------------------------------------------------

RED=$'\033[0;31m'; GREEN=$'\033[0;32m'; YELLOW=$'\033[1;33m'
CYAN=$'\033[0;36m'; BOLD=$'\033[1m'; RESET=$'\033[0m'

RPC=""
REST=""
SAMPLE_BLOCKS=100   # number of past blocks used to compute avg block time
BUFFER=500          # extra blocks added on top of the vote-end estimate
EXPEDITED=false

die()  { echo -e "${RED}Error:${RESET} $*" >&2; exit 1; }
info() { echo -e "${CYAN}$*${RESET}"; }
ok()   { echo -e "${GREEN}$*${RESET}"; }

usage() {
  cat <<EOF
${BOLD}Usage:${RESET} $(basename "$0") [flags]

  --rpc    <url>  RPC endpoint   e.g. http://localhost:26657  (required)
  --rest   <url>  REST endpoint  e.g. http://localhost:1317   (required)
  --sample <n>    Blocks to sample for avg block time         default: 100
  --buffer <n>    Extra blocks added after vote-end estimate  default: 500
  --expedited     Use expedited_voting_period instead of voting_period
  -h, --help      Show this help
EOF
  exit 0
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --rpc)    RPC="$2";          shift 2 ;;
    --rest)   REST="$2";         shift 2 ;;
    --sample)    SAMPLE_BLOCKS="$2"; shift 2 ;;
    --buffer)    BUFFER="$2";       shift 2 ;;
    --expedited) EXPEDITED=true;    shift   ;;
    -h|--help)   usage ;;
    *) die "Unknown flag: $1. Use --help for usage." ;;
  esac
done

[[ -z "$RPC"  ]] && die "--rpc is required"
[[ -z "$REST" ]] && die "--rest is required"

# ---------------------------------------------------------------------------
# Helper: fetch JSON and extract a field using python3
# ---------------------------------------------------------------------------
fetch() { curl -sf --max-time 10 "$1"; }

py() { python3 -c "$@"; }

# ---------------------------------------------------------------------------
# 1. Current block height + timestamp
# ---------------------------------------------------------------------------
info "Fetching latest block..."

LATEST_JSON=$(fetch "${REST}/cosmos/base/tendermint/v1beta1/blocks/latest") \
  || die "Cannot reach REST endpoint: ${REST}"

CURRENT_HEIGHT=$(py "import json,sys; d=$LATEST_JSON; print(int(d['block']['header']['height']))")
LATEST_TIME=$(py "import json,sys; d=$LATEST_JSON; print(d['block']['header']['time'])")

ok "Current height : ${BOLD}${CURRENT_HEIGHT}${RESET}"
ok "Latest block time : ${BOLD}${LATEST_TIME}${RESET}"

# ---------------------------------------------------------------------------
# 2. Block N ago + timestamp → average block time
# ---------------------------------------------------------------------------
SAMPLE_BLOCKS=$(( SAMPLE_BLOCKS > CURRENT_HEIGHT ? CURRENT_HEIGHT - 1 : SAMPLE_BLOCKS ))
OLD_HEIGHT=$(( CURRENT_HEIGHT - SAMPLE_BLOCKS ))

info "Fetching block ${OLD_HEIGHT} (${SAMPLE_BLOCKS} blocks ago) for avg block time..."

OLD_JSON=$(fetch "${REST}/cosmos/base/tendermint/v1beta1/blocks/${OLD_HEIGHT}") \
  || die "Cannot fetch block ${OLD_HEIGHT}"

OLD_TIME=$(py "import json; d=$OLD_JSON; print(d['block']['header']['time'])")

# Calculate avg block time in seconds using python datetime
AVG_BLOCK_TIME=$(py "
import json
from datetime import datetime, timezone

def parse(t):
    # Handle nanoseconds: truncate to microseconds
    t = t.rstrip('Z')
    if '.' in t:
        base, frac = t.split('.')
        frac = frac[:6].ljust(6, '0')
        t = base + '.' + frac
    return datetime.fromisoformat(t).replace(tzinfo=timezone.utc)

latest = parse('${LATEST_TIME}')
old    = parse('${OLD_TIME}')
diff   = (latest - old).total_seconds()
avg    = diff / ${SAMPLE_BLOCKS}
print(f'{avg:.3f}')
")

ok "Avg block time : ${BOLD}${AVG_BLOCK_TIME}s${RESET} (sampled over ${SAMPLE_BLOCKS} blocks)"

# ---------------------------------------------------------------------------
# 3. Governance params
# ---------------------------------------------------------------------------
info "Fetching governance params..."

GOV_JSON=$(fetch "${REST}/cosmos/gov/v1/params/voting") \
  || die "Cannot fetch gov params from ${REST}"

# Parse all relevant gov params in one python call.
# JSON is piped via stdin to avoid shell-expanding JSON null as a Python name.
# Response structure (SDK v0.53):
#   params.voting_period            — canonical location
#   params.max_deposit_period
#   params.min_deposit[]
#   params.expedited_voting_period
#   params.quorum / threshold
#   voting_params.voting_period     — legacy field, used as fallback
GOV_PARSED=$(echo "$GOV_JSON" | python3 -c "
import json, sys

def secs(raw):
    if not raw:
        return 0
    return int(float(raw.rstrip('s')))

def human(s):
    h = s // 3600; m = (s % 3600) // 60
    return f'{h}h {m}m ({s}s)'

def pct(v):
    return f'{float(v)*100:.1f}%' if v else 'n/a'

d = json.load(sys.stdin)
p = d.get('params') or {}

# voting period — prefer params.voting_period, fall back to legacy voting_params
vp_raw = p.get('voting_period') or (d.get('voting_params') or {}).get('voting_period', '')
if not vp_raw:
    raise SystemExit('voting_period not found in gov params response')

vp_secs = secs(vp_raw)

dep_raw  = p.get('max_deposit_period', '')
dep_secs = secs(dep_raw)

exp_raw  = p.get('expedited_voting_period', '')
exp_secs = secs(exp_raw)

min_dep = p.get('min_deposit') or []
min_dep_str = ', '.join(f\"{c.get('amount')} {c.get('denom')}\" for c in min_dep) if min_dep else 'n/a'

exp_dep = p.get('expedited_min_deposit') or []
exp_dep_str = ', '.join(f\"{c.get('amount')} {c.get('denom')}\" for c in exp_dep) if exp_dep else 'n/a'

print(vp_secs)
print(human(vp_secs))
print(dep_secs)
print(human(dep_secs))
print(exp_secs)
print(human(exp_secs))
print(min_dep_str)
print(exp_dep_str)
print(pct(p.get('quorum', '')))
print(pct(p.get('threshold', '')))
print(pct(p.get('veto_threshold', '')))
")

# Read each line into a variable
IFS=$'\n' read -r -d '' \
  VOTING_PERIOD_SECS \
  VOTING_PERIOD_HUMAN \
  DEPOSIT_PERIOD_SECS \
  DEPOSIT_PERIOD_HUMAN \
  EXPEDITED_PERIOD_SECS \
  EXPEDITED_PERIOD_HUMAN \
  MIN_DEPOSIT \
  EXPEDITED_MIN_DEPOSIT \
  QUORUM \
  THRESHOLD \
  VETO_THRESHOLD \
  <<< "$GOV_PARSED" || true

ok "Voting period  : ${BOLD}${VOTING_PERIOD_HUMAN}${RESET}"

# ---------------------------------------------------------------------------
# 4. Estimate block height when voting ends + buffer
# ---------------------------------------------------------------------------
if $EXPEDITED; then
  [[ "$EXPEDITED_PERIOD_SECS" == "0" ]] && die "expedited_voting_period is not set on this chain"
  ACTIVE_PERIOD_SECS="$EXPEDITED_PERIOD_SECS"
  ACTIVE_PERIOD_HUMAN="$EXPEDITED_PERIOD_HUMAN"
  ACTIVE_PERIOD_LABEL="Expedited voting period"
else
  ACTIVE_PERIOD_SECS="$VOTING_PERIOD_SECS"
  ACTIVE_PERIOD_HUMAN="$VOTING_PERIOD_HUMAN"
  ACTIVE_PERIOD_LABEL="Voting period"
fi

BLOCKS_IN_VOTE=$(py "import math; print(math.ceil(${ACTIVE_PERIOD_SECS} / ${AVG_BLOCK_TIME}))")
VOTE_END_HEIGHT=$(( CURRENT_HEIGHT + BLOCKS_IN_VOTE ))
UPGRADE_HEIGHT=$(( VOTE_END_HEIGHT + BUFFER ))

# Human-readable estimate of when upgrade height is reached
TOTAL_BLOCKS=$(( UPGRADE_HEIGHT - CURRENT_HEIGHT ))
TOTAL_SECS=$(py "print(int(${TOTAL_BLOCKS} * ${AVG_BLOCK_TIME}))")
UPGRADE_ETA=$(py "
from datetime import datetime, timedelta, timezone
eta = datetime.now(timezone.utc) + timedelta(seconds=${TOTAL_SECS})
print(eta.strftime('%Y-%m-%d %H:%M:%S UTC'))
")

VOTE_END_ETA=$(py "
from datetime import datetime, timedelta, timezone
eta = datetime.now(timezone.utc) + timedelta(seconds=int(${BLOCKS_IN_VOTE} * ${AVG_BLOCK_TIME}))
print(eta.strftime('%Y-%m-%d %H:%M:%S UTC'))
")

# ---------------------------------------------------------------------------
# 5. Print result
# ---------------------------------------------------------------------------
echo ""
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${BOLD}Governance Params${RESET}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
printf "  %-30s %s\n" "Min deposit:"              "${MIN_DEPOSIT}"
printf "  %-30s %s\n" "Deposit period:"           "${DEPOSIT_PERIOD_HUMAN}"
printf "  %-30s %s\n" "Voting period:"            "${VOTING_PERIOD_HUMAN}"
printf "  %-30s %s\n" "Expedited voting period:"  "${EXPEDITED_PERIOD_HUMAN}"
printf "  %-30s %s\n" "Expedited min deposit:"    "${EXPEDITED_MIN_DEPOSIT}"
printf "  %-30s %s\n" "Quorum:"                   "${QUORUM}"
printf "  %-30s %s\n" "Threshold:"                "${THRESHOLD}"
printf "  %-30s %s\n" "Veto threshold:"           "${VETO_THRESHOLD}"
echo ""
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${BOLD}Upgrade Height Estimate${RESET}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
printf "  %-30s %s\n" "Current height:"           "${CURRENT_HEIGHT}"
printf "  %-30s %s\n" "Avg block time:"           "${AVG_BLOCK_TIME}s"
printf "  %-30s %s\n" "${ACTIVE_PERIOD_LABEL}:"   "${ACTIVE_PERIOD_HUMAN}"
printf "  %-30s %s\n" "Blocks during vote:"       "${BLOCKS_IN_VOTE}"
printf "  %-30s %s\n" "Estimated vote-end height:" "${VOTE_END_HEIGHT}  (~${VOTE_END_ETA})"
printf "  %-30s %s\n" "Buffer blocks:"            "${BUFFER}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
printf "  %-30s ${GREEN}${BOLD}%s${RESET}\n" "Recommended upgrade height:" "${UPGRADE_HEIGHT}"
printf "  %-30s ${GREEN}${BOLD}%s${RESET}\n" "Estimated upgrade time:"     "${UPGRADE_ETA}"
echo -e "${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo ""
echo -e "${YELLOW}Use this height with set-upgrade-height.sh:${RESET}"
echo -e "  ${CYAN}./set-upgrade-height.sh --upgrade-name <name> --height ${UPGRADE_HEIGHT} --from <key> --submit${RESET}"
echo ""
