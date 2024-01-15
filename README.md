# goth-esa

[![CI](https://github.com/winebarrel/goth-esa/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/goth-esa/actions/workflows/ci.yml)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/winebarrel/goth-esa)](https://github.com/winebarrel/goth-esa/tags)

[Goth](https://github.com/markbates/goth) provider for [esa](https://docs.esa.io/posts/102#OAuth%E3%82%92%E5%88%A9%E7%94%A8%E3%81%97%E3%81%9F%E8%AA%8D%E5%8F%AF%E3%83%95%E3%83%AD%E3%83%BC).

## Getting Started

See [the official documentation](https://docs.esa.io/posts/102#OAuth%E3%82%92%E5%88%A9%E7%94%A8%E3%81%97%E3%81%9F%E8%AA%8D%E5%8F%AF%E3%83%95%E3%83%AD%E3%83%BC).

1. Register your application from `https://[team].esa.io/user/applications`
    * Redirect URI: `http://localhost:3000/auth/esa/callback`
1. Set the following environment variables
    * `ESA_KEY`: OAuth2 Client ID
    * `ESA_SECRET`: OAuth2 Secret
1. Run `make`
1. Open http://localhost:3000

## Usage

See https://github.com/winebarrel/goth-esa/blob/main/_example/main.go
