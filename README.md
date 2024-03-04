# livefetcher

IMPORTANT: This program depends on mecab, please see https://github.com/shogo82148/go-mecab for install instructions. livefetcher Makefile expects mecab-config to be installed, and will set up environment flags appropriately if you use it. **Note that if all you are doing is creating a new connector, you do not actually need to run the software, only run tests, which do not require mecab to be installed. See wiki for details.**

To run hot-reloaded basic dev version of livefetcher, simply run `make migrate` and then `make watch`. This will require you to have postgres installed on your machine and gow set up.

To run containerized, run `make run`

## Connector development

See wiki once I finish writing it.

The relevant code to connector development is also pretty heavily documented.
