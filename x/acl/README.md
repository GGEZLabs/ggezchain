# x/acl

## Abstract

The `x/acl` module provides an Access Control List for the ggezchain network.
The ACL module supports permission assignment through `AccessDefinition`
and enforces maker-checker principles where applicable.

## Contents

- [Abstract](#abstract)
- [State](#state)
  - [SuperAdmin](#superAdmin)
  - [Admin](#admin)
  - [AclAuthority](#aclAuthority)
- [Messages](#messages)
  - [MsgInit](#msgInit)
  - [MsgUpdateSuperAdmin](#msgUpdateSuperAdmin)
  - [MsgAddAdmin](#MsgAddAdmin)
  - [MsgDeleteAdmin](#MsgDeleteAdmin)
  - [MsgAddAuthority](#MsgAddAuthority)
  - [MsgUpdateAuthority](#MsgUpdateAuthority)
  - [MsgDeleteAuthority](#MsgDeleteAuthority)
- [Events](#events)
    - [Message Events](#message-events)
- [Client](#client)
    - [CLI](#cli)
    - [Query](#query)
    - [Transactions](#transactions)
---

## State

### SuperAdmin

A unique address that has the highest level of privilege. Only one exists at a time.

### Admin

A set of privileged accounts authorized to manage authorities.

### AclAuthority

Addresses with specific `AccessDefinition`s which determine allowed modules and permissions.

---

## Messages

In this section we describe the processing of the `acl` messages and the corresponding updates to the state. All created/modified state objects specified by each message are defined within the [state](#state) section.

### MsgInit

A super admin is initialized using the `MsgInit` message.

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L20
```

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L46-L50
```

This message is expected to fail if the super admin already initialized.

### MsgUpdateSuperAdmin

A super admin can be updated using the `MsgUpdateSuperAdmin` message.

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L21
```

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L54-L58
```

This message is expected to fail if:

* the super admin does not initialized.
* signer is not the super admin.

### MsgAddAdmin

One or more admins can be added using the `MsgAddAdmin` message.

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L22
```

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L62-L66
```

This message is expected to fail if:

* signer is not the super admin.
* admin already exists.

### MsgDeleteAdmin

One or more admins can be deleted using the `MsgDeleteAdmin` message.

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L23
```

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L70-L74
```

This message is expected to fail if:

* signer is not the super admin.
* admin does not exist.

### MsgAddAuthority

An authority can be added using the `MsgAddAuthority` message.

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L24
```

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L78-L84
```

This message is expected to fail if:

* the super admin does not initialized.
* the signer is not a super admin or admin.
* authority address already exists.
* invalid access definition list format.
* empty access definition list.
* add duplicate module names.

### MsgUpdateAuthority

An authority can be updated using the `MsgUpdateAuthority` message.

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L25
```

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L88-L98
```

This message is expected to fail if:

* the super admin does not initialized.
* the signer is not a super admin or admin.
* authority address does not exist.
* try to add/update and remove same module in the same request.

### MsgDeleteAuthority

An authority can be deleted using the `MsgDeleteAuthority` message.

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L26
```

```protobuf reference
https://github.com/GGEZLabs/ggezchain/blob/main/proto/ggezchain/acl/tx.proto#L102-L106
```

This message is expected to fail if:

* the super admin does not initialized.
* the signer is not a super admin or admin.
* authority address does not exist.

---
## Events

The `acl` module emits the following events:

### MsgInit

| Type     | Attribute Key | Attribute Value |
| -------- | ------------- | --------------- |
| init     | super_admin   | {superAdmin}    |

### MsgUpdateSuperAdmin

| Type               | Attribute Key | Attribute Value |
| ------------------ | ------------- | --------------- |
| update_super_admin | super_admin   | {newSuperAdmin} |

### MsgAddAdmin

| Type      | Attribute Key | Attribute Value |
| --------- | ------------- | --------------- |
| add_admin | admins        | {admins}        |

### MsgDeleteAdmin

| Type         | Attribute Key | Attribute Value |
| ------------ | ------------- | --------------- |
| delete_admin | admins        | {admins}        |

### MsgAddAuthority

| Type          | Attribute Key      | Attribute Value     |
| ------------- | ------------------ | ------------------- |
| add_authority | authority_address  | {address}           |
| add_authority | name               | {name}              |
| add_authority | access_definitions | {accessDefinitions} |

### MsgUpdateAuthority

| Type             | Attribute Key      | Attribute Value     |
| ---------------- | ------------------ | ------------------  |
| update_authority | authority_address  | {address}           |
| update_authority | name               | {name}              |
| update_authority | access_definitions | {accessDefinitions} |

### MsgDeleteAuthority

| Type             | Attribute Key     | Attribute Value |
| ---------------- | ----------------- | --------------- |
| delete_authority | authority_address | {address}       |

---

## Client

### CLI

A user can query and interact with the `acl` module using the CLI.

#### Query

The `query` commands allow users to query `acl` state.

```shell
ggezchaind query acl --help
```

##### super-admin

The `super-admin` command allows users to query super-admin.

```shell
ggezchaind query acl super-admin [flags]
```

Example:

```shell
ggezchaind query acl super-admin
```

Example Output:

```yml
super_admin:
  admin: ggez1..
```

##### admins

The `admins` command allows users to query all admins.

```shell
ggezchaind query acl admins [flags]
```

Example:

```shell
ggezchaind query acl admins
```

Example Output:

```yml
acl_admin:
- address: ggez1q3sfaepes35ly4sa5ppguf6gs49un4uz858e33
pagination:
  total: "1"
```

##### admin

The `admin` command allows users to query specified admin.

```shell
ggezchaind query acl admin [address] [flags]
```

Example:

```shell
ggezchaind query acl admin ggez1..
```

Example Output:

```yml
acl_admin:
  address: ggez1q3sfaepes35ly4sa5ppguf6gs49un4uz858e33
```

##### acl-authorities

The `acl-authorities` command allows users to query all acl-authorities.

```shell
ggezchaind query acl acl-authorities [flags]
```

Example:

```shell
ggezchaind query acl acl-authorities
```

Example Output:

```yml
acl_authority:
- access_definitions:
  - is_checker: false
    is_maker: true
    module: trade
  address: ggez1..
  name: Alice
pagination:
  total: "1"
```

##### acl-authority

The `acl-authority` command allows users to query all acl-authority.

```shell
ggezchaind query acl acl-authority [address] [flags]
```

Example:

```shell
ggezchaind query acl acl-authority ggez1..
```

Example Output:

```yml
acl_authority:
  access_definitions:
  - is_checker: false
    is_maker: true
    module: trade
  address: ggez1..
  name: Alice
```

---

#### Transactions

The `tx` commands allow users to interact with the `acl` module.

```shell
ggezchaind tx acl --help
```

##### init

The `init` command initializes the super-admin. Can only be called once.

```shell
ggezchaind tx acl init [super-admin] [flags]
```

Example:

```shell
ggezchaind tx acl init ggez1..
```

##### update-super-admin

The `update-super-admin` command update super admin. Only a super admin can perform this action.

```shell
ggezchaind tx acl update-super-admin [new-super-admin] [flags]
```

Example:

```shell
ggezchaind tx acl update-super-admin ggez1..
```

##### add-admin

The `add-admin` command add one or more admin. Only a super admin can perform this action.

```shell
ggezchaind tx acl add-admin [admins] [flags]
```

Example:

```shell
ggezchaind tx acl add-admin ggez1..,ggez1..
```

##### delete-admin

The `delete-admin` command delete one or more admin. Only a super admin can perform this action.

```shell
ggezchaind tx acl delete-admin [admins] [flags]
```

Example:

```shell
ggezchaind tx acl delete-admin ggez1..,ggez1..
```

##### add-authority

The `add-authority` command add a new authority with specific access definition. Must have authority to do so.

```shell
ggezchaind tx acl add-authority [auth-address] [name] [access-definitions] [flags]
```

Example:

```shell
ggezchaind tx acl add-authority ggez1.. Alice '[{"module":"trade","is_maker":true,"is_checker":false}]'
```

##### add-authority

The `update-authority` command modify the name or access definition of an existing authority. Must have authority to do so.

```shell
ggezchaind tx acl update-authority [auth-address] [name] [access-definitions] [flags]
```

Example:

```shell
ggezchaind tx acl update-authority ggez1.. Alice '[{"module":"trade","is_maker":true,"is_checker":false}]'
```

