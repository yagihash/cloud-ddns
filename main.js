const config = require('config').config;
const {DNS} = require('@google-cloud/dns');
const rp = require('request-promise');

const fetchCurrentIp = async () => {
  const { ip: currentIp } = await rp({
    uri: 'https://ipinfo.io/',
    method: 'GET',
    json: true,
  });
  return currentIp;
};

const fetchActiveIp = async (zone, domainName) => {
  const [[{ data: [activeIp] }]] = await zone.getRecords({
    name: domainName,
    type: 'A',
  });
  return activeIp;
};

const createARecord = (zone, domainName, ttl, ip) => {
  return zone.record('a', {
    name: domainName,
    ttl: ttl,
    data: ip,
  });
};

const createRecordDiff = (zone, domainName, ttl, currentIp, activeIp) => {
  return {
    add: createARecord(zone, domainName, ttl, currentIp),
    delete: createARecord(zone, domainName, ttl, activeIp),
  }
};

const createSlackAttachments = (domainName, activeIp, currentIp) => JSON.stringify([{
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
  ts: Math.floor((new Date()).getTime() / 1000),
}]);

const postResultOnSlack = async (token, channel, username, attachments) => {
  return await rp({
    uri: 'https://slack.com/api/chat.postMessage',
    method: 'POST',
    form: {
      token: token,
      channel: channel,
      username: username,
      attachments: attachments,
    },
    json: true,
  });
};

const main = async () => {
  const projectId = config.projectId;
  const managedZones = config.managedZones;

  const dns = new DNS({
    projectId: projectId,
  });

  const currentIp = await fetchCurrentIp();

  for (let managedZone of managedZones) {
    let { name: zoneName, ttl } = managedZone;
    for (let domainName of managedZone.domainNames) {
      const zone = dns.zone(zoneName);
      const activeIp = await fetchActiveIp(zone, domainName);
      if (currentIp !== activeIp) {
        const diff = createRecordDiff(zone, domainName, ttl, currentIp, activeIp);

        zone.createChange(diff).then(() => {
          (async () => {
            const slackConfig = config.slack;
            const { token, channel, username } = slackConfig;
            const attachments = createSlackAttachments(domainName, activeIp, currentIp);
            await postResultOnSlack(token, channel, username, attachments);
          })();
        });
      }
    }
  }
};

main().catch((e) => console.log(e));
