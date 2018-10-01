# slack-new-channel

## Usage

1. Download latest `slack-new-channel` binary that suits your platform from [release page ](https://github.com/kitsuyui/slack-new-channel/releases/).
2. Make it executable. (`chmod +x ./slack-new-channel`)
3. Set environment variables.
4. Execute it.

```
$ ./slack-new-channel
```

And then it reports new channels into incoming-webhooks channel.
<br>
<img src="https://user-images.githubusercontent.com/2596972/46255289-35d50b00-c4d6-11e8-88c9-56e01f0053ce.png" width="50%" height="50%" title="Example">

## Run periodically

It doesn't have any daemonizing options.
So you must set in cron jobs or daemonize script if you want it runs periodically.

### Tokens

- incoming-webhook: https://{ your organization }.slack.com/services/new/incoming-webhook
- oauth-tokens: https://api.slack.com/docs/oauth-test-tokens

### Environment variables.

See .env.sample .

```console
$ export LATEST_CHANNEL_JSON_PATH='./latest-created-channel.json'
$ export SLACK_API_TOKEN='xxxx-1234567890-1234567890-124567890-aaaaaaaaaaaaaaaaaaaaaaaa'
$ export SLACK_WEBHOOK_URL='https://hooks.slack.com/services/dUmMy/YOuR/HoOkUrLHeRE'
```

## Build

```
$ docker run --rm -v "$(pwd)":/slack-new-channel -w /slack-new-channel tcnksm/gox sh -c "./build.sh"
```
