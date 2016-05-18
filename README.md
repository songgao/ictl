# Interframe Compression Transmission Layer

Interframe Compression Transmission Layer (ICTL) compresses difference between consecutive messages rather than message themselves. For periodical messages where consecutive messages are similar to each other, the difference is normally fairly compressible, which results in smaller bandwidth consumption.

ICTL is designed in the context of vehicular networking. However, it is implemented as a generic lossless compression layer that can be useful for transmsision of any periodical binary data over unreliable transport.

In some tests, we have achieved over 50% bandwidth saving in vehicular binary data.

Please see [https://song.gao.io/dissertation/](https://song.gao.io/dissertation/) for more details!


## License

[BSD 3-Clause License](./LICENSE)
