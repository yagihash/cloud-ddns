const config = require('config').config;
const DNS = require('@google-cloud/dns');
const rp = require('request-promise');

const projectId = config.projectId;
const managedZones = config.managedZones;

const dns = new DNS({
  projectId: projectId,
});

const zone = dns.zone('xss-moe');

// Old
rp({
  uri: 'https://ipinfo.io/',
  method: 'GET',
  json: true,
}).then(function (data) {
  const currentIp = data.ip;

  for (let managedZone of managedZones) {
    for (let domainName of managedZone.domainNames) {
      zone.getRecords({
        name: domainName,
        type: 'A',
      }).then(function (data) {
        const activeIp = data[0][0].data[0];
        if (activeIp !== currentIp) {
          const newRecord = zone.record('a', {
            name: domainName,
            ttl: managedZone.ttl,
            data: currentIp,
          });

          const oldRecord = zone.record('a', {
            name: domainName,
            ttl: managedZone.ttl,
            data: activeIp,
          });

          const diff = {
            add: newRecord,
            delete: oldRecord,
          }

          zone.createChange(diff).then(function (data) {
            const token = config.slack.token;
            const channel = config.slack.channel;
            const username = config.slack.username;
            const iconEmoji = config.slack.iconEmoji;

            rp({
              uri: 'https://slack.com/api/chat.postMessage',
              method: 'POST',
              form: {
                token: token,
                channel: channel,
                username: username,
                icon_emoji: iconEmoji,
                attachments: JSON.stringify([{
                  pretext: 'Cloud-DDNS updated your record.',
                  fallback: `Updated A record for ${domainName} from ${activeIp} to ${currentIp}`,
                  color: 'good',
                  title: domainName,
                  title_link: `http://${domainName.slice(-1)}/`,
                  fields: [
                    {
                      title: 'Previous record',
                      value: activeIp,
                      short: true,
                    },
                    {
                      title: 'New record',
                      value: currentIp,
                      short: true,
                    },
                  ],
                  footer: 'Sent via cloud-ddns',
                  footer_icon: ':gcp:',
                  ts: Math.floor((new Date()).getTime() / 1000),
                }]),
              },
              json: true,
            });
          });
        }
      });
    }
  }
});

// With async-await
const proposal = async () => {
  const data = await rp({
    uri: 'https://ipinfo.io/',
    method: 'GET',
    json: true,
  }); // const { ip } = ... is more readable and useful

  const currentIp = data.ip;

  for (let managedZone of managedZones) {
    for (let domainName of managedZone.domainNames) {

      const data = await zone.getRecords({
        name: domainName,
        type: 'A',
      });

      const activeIp = data[0][0].data[0];
      if (activeIp !== currentIp) {
        const newRecord = zone.record('a', {
          name: domainName,
          ttl : managedZone.ttl,
          data: currentIp,
        });

        const oldRecord = zone.record('a', {
          name: domainName,
          ttl : managedZone.ttl,
          data: activeIp,
        });

        const diff = {
          add   : newRecord,
          delete: oldRecord,
        };

        await zone.createChange(diff);

        const token = config.slack.token;
        const channel = config.slack.channel;
        const username = config.slack.username;
        const iconEmoji = config.slack.iconEmoji;

        await rp({
          uri   : 'https://slack.com/api/chat.postMessage',
          method: 'POST',
          form  : {
            token      : token,
            channel    : channel,
            username   : username,
            icon_emoji : iconEmoji,
            attachments: JSON.stringify([{
              pretext    : 'Cloud-DDNS updated your record.',
              fallback   : `Updated A record for ${domainName} from ${activeIp} to ${currentIp}`,
              color      : 'good',
              title      : domainName,
              title_link : `http://${domainName.slice(-1)}/`,
              fields     : [
                {
                  title: 'Previous record',
                  value: activeIp,
                  short: true,
                },
                {
                  title: 'New record',
                  value: currentIp,
                  short: true,
                },
              ],
              footer     : 'Sent via cloud-ddns',
              footer_icon: ':gcp:',
              ts         : Math.floor((new Date()).getTime() / 1000),
            }]),
          },
          json  : true,
        });
      }
    }
  }
};

proposal()
  .then(() => { console.log('success'); })
  .catch(e => { throw e; });

// More simplfy
const proposal2 = async () => {
  const { ip: currentIp } = await rp({
    uri: 'https://ipinfo.io/',
    method: 'GET',
    json: true,
  });

  const getRecordsPromises = [];
  const tmpDomainNames = [];
  const tmpManagedZones = [];

  for (let managedZone of managedZones) {
    for (let domainName of managedZone.domainNames) {
      getRecordsPromises.append(zone.getRecords({
        name: domainName,
        type: 'A',
      }));
      tmpDomainNames.append(domainName);
      tmpManagedZones.append(managedZone);
    }
  }

  const getRecordsResults = await Promise.all(getRecordsPromises);

  for (let i = 0;i < getRecordsResults.length; i++) {
    const result = getRecordsResults[i];
    const domainName = tmpDomainNames[i];
    const managedZone = tmpManagedZones[i];

    const activeIp = result[0][0].data[0];

    if (activeIp !== currentIp) {
      const newRecord = zone.record('a', {
        name: domainName,
        ttl : managedZone.ttl,
        data: currentIp,
      });

      const oldRecord = zone.record('a', {
        name: domainName,
        ttl : managedZone.ttl,
        data: activeIp,
      });

      const diff = {
        add   : newRecord,
        delete: oldRecord,
      };

      await zone.createChange(diff);

      const slackConfig = config.slack;
      const { token, channel, username, iconEmoji } = slackConfig;

      await rp({
        uri   : 'https://slack.com/api/chat.postMessage',
        method: 'POST',
        form  : {
          token      : token,
          channel    : channel,
          username   : username,
          icon_emoji : iconEmoji,
          attachments: JSON.stringify([{
            pretext    : 'Cloud-DDNS updated your record.',
            fallback   : `Updated A record for ${domainName} from ${activeIp} to ${currentIp}`,
            color      : 'good',
            title      : domainName,
            title_link : `http://${domainName.slice(-1)}/`,
            fields     : [
              {
                title: 'Previous record',
                value: activeIp,
                short: true,
              },
              {
                title: 'New record',
                value: currentIp,
                short: true,
              },
            ],
            footer     : 'Sent via cloud-ddns',
            footer_icon: ':gcp:',
            ts         : Math.floor((new Date()).getTime() / 1000),
          }]),
        },
        json  : true,
      });
    }
  }
};

proposal2()
  .then(() => { console.log('success'); })
  .catch(e => { throw e; });
