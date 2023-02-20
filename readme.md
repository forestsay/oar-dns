# 椆-NS

Oar-NS (椆-NS) is a private DNS server used for resolving OpenVPN client name to OpenVPN Private Address.

Mount Buzhou (不周山) is one of the famous mountains that appear in Chinese mythology. In some literature, Mount Buzhou
is
a gap between living and death. We got the character "周" from the name of this mountain.

Wood (木) is one of the symbols of [ForestSay](https://github.com/forestsay). So we got the character "木" from it.

In glyph, "椆" is a combination of "木周". In meaning, "椆" is oar. Everybody in a developing group is just like
everyone on
the same boat. We need a tool to let us work more effectively, so maybe we need an oar to help with.

## Build

1. Run `docker build .`

## Run

### Environment Variable

| Name                                    | Description                           | Default           |
|-----------------------------------------|---------------------------------------|-------------------|
| `OPENVPN-MANAGEMENT-INTERFACE-ENDPOINT` | Openvpn Management Interface Endpoint | `127.0.0.1:27273` |
| `OPENVPN-MANAGEMENT-INTERFACE-PASSWORD` | Openvpn Management Interface Password | `password`        |
| `DNS-SERVER-LISTEN`                     | Dns server listen address             | `:53`             |

