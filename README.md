# TSID Generator

Go library for generating Time-Sorted Unique Identifiers (TSID)
Implementation inspired from [f4b6a3/tsid-creator](https://github.com/f4b6a3/tsid-creator)

It brings together ideas from [Twitter's Snowflake](https://github.com/twitter-archive/snowflake/tree/snowflake-2010) and [ULID Spec](https://github.com/ulid/spec).

In summary:

- Sorted by generation time;
- Can be stored as an integer of 64 bits;
- Can be stored as a string of 13 chars;
- String format is encoded to [Crockford's base32](https://www.crockford.com/base32.html);
- String format is URL safe, is case insensitive, and has no hyphens;
- Shorter than UUID, ULID and KSUID.

Recommended readings:

- [The best UUID type for a database Primary Key](https://vladmihalcea.com/uuid-database-primary-key/)
- [The primary key dilemma: ID vs UUID and some practical solutions](https://fillumina.wordpress.com/2023/02/06/the-primary-key-dilemma-id-vs-uuid-and-some-practical-solutions/)
- [Primary keys in the DB - what to use? ID vs UUID or is there something else?](https://www.linkedin.com/pulse/primary-keys-db-what-use-id-vs-uuid-something-else-lucas-persson)

## Usage

Create a TSID Factory:

```go
tsidFactory, err := TsidFactoryBuilder().
    WithNodeBits(nodeBits). // max 20
    WithNode(nodeId). // max 2^nodeBits
    WithCustomEpoch(customEpoch).
    WithClock(clock).
    WithRandom(random).
    Build()
```

> Default node Id is 0 & node bits is 0. Epoch and clock are optional, and Random when not provided it internally uses
> random from math/rand to generate random values

> [!NOTE]
> `.Build()` creates / returns the existing instance of the tsid factory, this is useful where single
> instance of tsid factory needs to be shared across go routines. But if you need a new instance
> use `.NewInstance()`

Generate TSID

```go
tsid, err := tsidFactory.Generate()
```

Get TSID as `int64`:

```go
number := tsid.ToNumber(); // 122390644586507544
```

Get TSID as `string`:

```go
tsidStr := tsid.ToString(); // 03CPHMJ76HV8R
```

The TSID generator is [thread-safe](https://en.wikipedia.org/wiki/Thread_safety).

### Dependency

Run the following command:

```shell
go get github.com/rushysloth/go-tsid
```

### TSID as Number

The `tsid.ToNumber()` method simply unwraps the internal `int64` value of a TSID.

```go
id := tsid.ToNumber()
```

Sequence of TSIDs:

```text
38352658567418867
38352658567418868
38352658567418869
38352658567418870
38352658567418871
38352658567418872
38352658567418873
38352658567418874
38352658573940759 < millisecond changed
38352658573940760
38352658573940761
38352658573940762
38352658573940763
38352658573940764
38352658573940765
38352658573940766
         ^      ^ look

|--------|------|
   time   random
```

### TSID as String

The `tsid.ToString()` method encodes a TSID to [Crockford's base 32](https://www.crockford.com/base32.html) encoding. The returned string is 13 characters long.

```go
idStr := tsid.ToString();
```

Sequence of TSID strings:

```text
01226N0640J7K
01226N0640J7M
01226N0640J7N
01226N0640J7P
01226N0640J7Q
01226N0640J7R
01226N0640J7S
01226N0640J7T
01226N0693HDA < millisecond changed
01226N0693HDB
01226N0693HDC
01226N0693HDD
01226N0693HDE
01226N0693HDF
01226N0693HDG
01226N0693HDH
        ^   ^ look

|-------|---|
   time random
```

The string format can be useful for languages that store numbers in [double-precision 64-bit binary format IEEE 754](https://en.wikipedia.org/wiki/Double-precision_floating-point_format), such as [Javascript](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number).

### TSID Structure

The term TSID stands for (roughly) Time-Sorted ID. A TSID is a number that is formed by the creation time along with random bits.

The TSID has 2 components:

- Time component (42 bits)
- Random component (22 bits)

The time component is the count of milliseconds since 2020-01-01 00:00:00 UTC.

The Random component has 2 sub-parts:

- Node ID (0 to 20 bits)
- Counter (2 to 22 bits)

The counter bits depend on the node bits. If the node bits are 10, the counter bits are limited to 12. In this example, the maximum node value is 2^10-1 = 1023 and the maximum counter value is 2^12-1 = 4095. So the maximum TSIDs that can be generated per millisecond is 4096.

The node identifier uses 10 bits of the random component by default in the `TsidFactory`. It's possible to adjust the node bits to a value between 0 and 20. The counter bits are affected by the node bits.

This is the default TSID structure:

```
                                            adjustable
                                           <---------->
|------------------------------------------|----------|------------|
       time (msecs since 2020-01-01)           node       counter
                42 bits                       10 bits     12 bits

- time:    2^42 = ~69 years or ~139 years (with adjustable epoch)
- node:    2^10 = 1,024 (with adjustable bits)
- counter: 2^12 = 4,096 (initially random)

Notes:
The node is adjustable from 0 to 20 bits.
The node bits affect the counter bits.
The time component can be used for ~69 years if stored in a SIGNED 64 bits integer field.
The time component can be used for ~139 years if stored in a UNSIGNED 64 bits integer field.
```

The time component can be 1 ms or more ahead of the system time when necessary to maintain monotonicity and generation speed.

### Node identifier

A simple way to avoid collisions is to make sure that each generator has its exclusive node ID. A "node" as we call it in this library can be a physical machine, a virtual machine, a container, a k8s pod, a running process, a database instance number, etc.

**Notes:**

1. As a reference, [6,000 tweets are posted on Twitter every second as of 2022](https://www.demandsage.com/twitter-statistics/);
1. According to the `hostname` manual:
   - `--ip-address` or `-i` (lowercase): _Display the network address(es) of the host name. Note that this works only if the host name can be resolved. Avoid using this option; use hostname --all-ip-addresses instead_.
   - `--all-ip-addresses` or `-I` (uppercase): _Display all network addresses of the host. This option enumerates all configured addresses on all network interfaces. The loopback interface and IPv6 link-local addresses are omitted. Contrary to option -i, this option does not depend on name resolution. Do not make any assumptions about the order of the output._

### More Examples

Create a quick TSID:

```go
tsid := tsid.Fast();
```

---

Create a quick TSID from canonical string (13 chars)

```go
tsid := tsid.FromString("03CPHMJ76HV8R")
```

---

Get the creation unix millis of the tsid

```go
millis := tsid.GetUnixMillis()
```

---

A `TsidFactory` with a FIXED node identifier and CUSTOM node bits:

```go
// setup a factory for up to 64 nodes and 65536 ID/ms.
tsidFactory, err := TsidFactoryBuilder()
    .WithNodeBits(6)      // max: 20
    .WithNode(63)         // max: 2^nodeBits
    .Build();

// use the factory
tsid, err := tsidFactory.Generate()
```

---

A `TsidFactory` with a CUSTOM epoch:

```go
epoch := time.Now().UnixMilli()
tsidFactory, err := TsidFactoryBuilder().
    WithCustomEpoch(epoch).
    build();

// use the factory
tsid, err := tsidFactory.Generate()
```

---

A `TsidFactory` with Crypto Random value generator

```go
// using crypto/rand for security use cases
supplier := NewCryptoRandomSupplier()

// creating random value generator
intRandom := NewIntRandom(supplier)

tsidFactory, err := TsidFactoryBuilder().
    WithRandom(intRandom).
    build();

```

> You can use custom random value suppliers either by implementing
> `IntSupplier`, `ByteSupplier` or using `NewIntRandomWithSupplierFunc`

---

## Ports, forks and other OSS

Ports and forks:

| Language | Name                                                                                  |
| -------- | ------------------------------------------------------------------------------------- |
| Java     | [vladmihalcea/hypersistence-tsid](https://github.com/vladmihalcea/hypersistence-tsid) |
| .NET     | [kgkoutis/TSID.Creator.NET](https://github.com/kgkoutis/TSID.Creator.NET)             |
| PHP      | [odan/tsid](https://github.com/odan/tsid)                                             |
| Python.  | [luismedel/tsid-python](https://github.com/luismedel/tsid-python)                     |

Other OSS:

| Language | Name                                                                |
| -------- | ------------------------------------------------------------------- |
| Java     | [fillumina/id-encryptor](https://github.com/fillumina/id-encryptor) |
| .NET     | [ullmark/hashids.net](https://github.com/ullmark/hashids.net)       |

## License

This library is Open Source software released under the [Apache-2.0 license](https://www.apache.org/licenses/LICENSE-2.0).
