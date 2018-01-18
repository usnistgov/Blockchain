# Model 1.0

Coin creation uses the "Trillion Dollar Platinum Coin" concept of
creating a single large value coin (here, 1 million units), and
spending it into existence incrementally through regular UBI
exercises. In version 0, CreateCoins used unit coins, and 1 million
outputs of denomination 1 require 1 million discrete outputs. Here,
the million unit coin is accomplished in a single line of output.

```
[471] Ledger:
    Tid: 1487547740451
    Ttyp: CreateCoins
    Id/Seq: 0/0 Denom: 0 Owner: MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQ
    CrtJSkjl+P06RUIACreJe7RTL5kBKUbYR+KNrvCbmjyKhvHbSWtLNAqrIHKyrcL0re
    eXfFS6IzFSxHsbr00b1t315IFxVuVADfnavZicTezMFyEapCDPDvJLxXeSPoXJ6hyx
    OB8SgeJ9Ted3HhNss4iUjs1pC1W85KKPuePFK4owIDAQAB
    Id/Seq: 1487547740451/0 Denom: 1000000 Owner: QSoT0u9VlnrL4wq2pppy
    +jB4lEbYJ7xeWnE1VKzgVic=
    Tsig:  N7zvw+ZyGzPlVCRy5cufve+UVhm5K+5PIn+zBH64iNfc7C4V9G9yYXiCLav
    fSmsCKKxFDspRGibolMvu8SwVyg5Bml5+ZaNlpBejgzJSNDcHd7kFuxUarOIQE5pWo
    6DJE1cHpPHcZdm0gafEdNXUWs3hJ3Ug7wkxVqKnqf4GBcY=
```

Large denomination coins achieve impressive orders of transaction size
compression, since no transaction requires more than 1 input and 2
outputs. This compares favorably against version 0, which for the
above requires 1 million inputs and 1 million outputs. This explains
why, not only is the domestic currency divided into larger and smaller
bills and coins, but also the checking system allows a single voucher
to transact arbitrarily large amounts of money.

Corollary to this improvement is that most transactions will now 'give
change' back to the payer. Much as you get $19 change back when you
give a $20 bill for a $1 candy bar, a payer who has a 20 unit coin to
hand over in payment for a 1 unit obligation, puts the 20 unit coin
as a single input, and generates 2 outputs: a 1 unit coin to the payee,
and a 19 unit coin back to self. Inputs and Outputs must balance, and
the chain of ownership is traceable, since the coin owner is identified
by public key and the transaction outputs are signed by the corresponding
private key. Here is an instance of the UBI payment with the banker starting
at a 1010000 coin, distributing 20 to each recipient and returning
1009840 change to self.

```
[4] Ledger:
    Tid: 1487946650275
    Ttyp: PayCoins

Id/Seq: 14879270/0 Denom: 1010000 Owner: QSoT0

Id/Seq: 14879275/0 Denom: 20 Owner: wVWP7p
Id/Seq: 14879275/1 Denom: 20 Owner: XLuz8U
Id/Seq: 14879275/2 Denom: 20 Owner: 2MwzAJ
Id/Seq: 14879275/3 Denom: 20 Owner: /4fvUk
Id/Seq: 14879275/4 Denom: 20 Owner: tGjSEm
Id/Seq: 14879275/5 Denom: 20 Owner: vK1HKs
Id/Seq: 14879275/6 Denom: 20 Owner: uQREkQ
Id/Seq: 14879275/7 Denom: 20 Owner: jlvgVr
Id/Seq: 14879275/8 Denom: 1009840 Owner: QSoT0u

    Tsig:  VWrR5V ... g6GUQNStO5KiEcq5NYfTTZmc=
```

Components remain the same as for Model 0.
