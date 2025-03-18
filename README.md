# livefetcher

IMPORTANT: This program depends on mecab, please see https://github.com/shogo82148/go-mecab for install instructions. livefetcher Makefile expects mecab-config to be installed, and will set up environment flags appropriately if you use it. **Note that if all you are doing is creating a new connector, you do not actually need to run the software, only run tests, which do not require mecab to be installed. See wiki for details.**

To run hot-reloaded basic dev version of livefetcher, simply run `make migrate` and then `make watch`. This will require you to have postgres installed on your machine and gow set up.

To run containerized, run `make run`

## Connector development

See wiki.

The relevant code to connector development is also pretty heavily documented.

Check [livehouse todos](LIVEHOUSE-TODO.md) for some livehouses to look into being implemented

## Release roadmap
### Immediate plans
- [x] add coordinates to live houses (these must be added to every existing live house too, it would be good to get this done before backlog is too large)
- [x] create design documents
- [ ] complete same live recognition for update functionality (currently about 95%+ accurate, but should be very close to 100% as many functions would rely on it)

### Relatively near future
- [x] add user auth system (required for many future features)
- [x] implement proper frontend
- [x] pagination of SELECT query
- [ ] batch live/livelist favorite SELECT queries

### Search improvements
- [x] implement search within distance from point
- [x] implement saved searches with notification functionality
- [x] add live to calendar (google calendar, maybe more)
- [x] bookmark lives
- [x] bookmarked lives notifications
- [x] add live lists

### Artist digest
- [ ] improve artist recognition to make artist names more uniform (this will likely imply creating two columns for artists - one with any extra info, and one without)
- [ ] create artist info repo
- [ ] create github bot to commit changes to artists

### Independent, but necessary
- [ ] implement 300 live houses in tokyo (primarily shimokitazawa, shibuya, shinjuku)
- [ ] improve docker containers, improve database resiliency

### If possible
- [ ] refactor the core router
- [ ] improve logging (It would probably be good to improve the core router in order to do this)

### Post-release
- [ ] create recommendation engine based on user data
