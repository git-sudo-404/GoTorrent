### GoTorrent Client

- A bittorent client implementation from scratch in go lang.

#### Here' How the repo is structured

```
.
├── cmd
│   └── main.go --> entrypoint for the GoTorrent
├── go.mod
├── internal
│   ├── bencode
│   │   ├── decoder_test.go
│   │   ├── decoder.go
│   │   ├── encoder_test.go
│   │   └── encoder.go
│   ├── torrent
│   │   ├── client.go
│   │   ├── messages.go
│   │   ├── metainfo.go
│   │   ├── metainfoParser_test.go
│   │   ├── metainfoParser.go
│   │   ├── torrent.go
│   │   ├── trackerRequest_test.go
│   │   └── trackerRequest.go
│   └── url-encoder
│       ├── urlencoder_test.go
│       └── urlencoder.go
├── LICENSE
├── README.md
└── testdata
    └── ubuntu-26.04.1-desktop-amd64.iso.torrent --> Test .torrent file
```

#### How to run

- Execute the following command after downloading a .torrent file and placing it under the testdata

```
go run cmd/main.go testdata/<your_torrent_file>.torrent
```

#### Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.
This project is licensed under the MIT License.
