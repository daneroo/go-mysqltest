# Consolidating restored snapshots December 2020

Given all the history of implementations and hardware adventures and failures,

- Scraping of the footprints software (Windows) (Aria Windows Media Server - 2008)
- Remote mounting the footprints data to linux, and doing ETL from SQLite -> MYSQL
- Implementing the python native scraper on linux
- Disk failures on goedel and later cantor
- Porting the scraper to Raspberry Pi (Which had a catastrophic Flash Card Failure) - that lasted a month 8-(
- Moving back to euler (née cantor) linux server (We lost data from 2016-02-14 to 2016-03-12)
- Re-implementing in Go (2018)
- Implementing the pump in Go for persistence neutral backups

The final rollup archive: `ted.20201120.2332Z.rollup-clean.jsonl.tar.bz2`, with IPFS/CID:`QmSLJPEZocdPZ99pazEkiJTaf3B1zeBmAQWEr7n9fSNgEu`

```bash
(cd data/jsonl; time tar cjvf ../../ted.20201120.2332Z.rollup-clean.jsonl.tar.bz2 ./)
```

## Systematic restore of all snapshots

- Older (legacy) Footprints sqlite databases were included in older snapshots (pre 2015-09-28)
- `ted.200YYMMDD.HHMM.sql.bz2` snapshots contain `ted_native`, `ted_service` and `watt` tables
  - we can ignore ted_service
    - seems in error, probably a time shift, we can ignore
    - has a problem with 2008-11-14T23:18:13Z - 2008-11-28T04:59:59Z
    - was in use only for [2008-11-14 23:18:13 , 2008-12-17 19:19:20] 
  - we can ignore `watt_day|hour|minute|tensec`

The tables `ted_native` and `watt`, should be equivalent starting at 2008-12-17T19:37:16Z, when the TedNative capture started. We only have dumps of both tables until: `ted.20150928.1006.sql.bz2`, and the last stamp is `2015-09-28T14:06:52Z`, at that point there are only 25 entries present in ted_native and missing in watt:

```txt
2020-12-11T05:41:34.294Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
```

So we combine both in a jsonl dump: `./data/jsonl-ted-rollup.20150928.1006/`, and confirm by restoring each snapshot (ted.200YYMMDD.HHMM.sql.bz2) and confirm that all data is in that combined json.

Phase-1 Verification: We then restored each `ted.200YYMMDD.HHMM.sql.bz2` to verify that all samples of `ted_native` and `watt` tables were included in the rollup. The output is in `RESTORE-phase-1.md`, (which took	24 hours), confirming that there were only `MissingInB`, that is no Conflict or MissingInA entries.

### Phase 2 - Final rollup - we will accumulate in postgres.
- Pre phase-2
  - seed with ./data/jsonl-ted-rollup.20150928.100 (result of phase-1 above)
  - restore these over postgres
    -	ted.watt.2016-02-14-1555.sql.bz2 - last watt backup including history from 2008-07-30 00:04:40
    - ted.watt-just2016.2016-02-14-1624.sql.gz - first backup with table truncated from 2016-01-01
    - ted.watt.20201120.2332Z.sql.bz2 - most recent backup with table truncated from 2016-01-01
- Phase-2 Verification
  - restore each `ted.watt.200YYMMDD.HHMM.sql.bz2` and verify against postgres.

### Cleanup

There was a problem, starting at `ted.watt.20181024.1913Z.sql.bz2` we have extra data:
```bash
2020-12-13T09:10:33.449Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInB
```

This data was removed from final rollup, as it's dates precede the purchase of the TED capture device, and it's data (< 1 day) could not be matched.

```sql
-- should be 75026 records
select count(*) from watt where stamp<'2007-08-29'
delete from watt where stamp<'2007-08-29'
```

### setup
```bash
mkdir -p ./data/archive/mirror/ted/
rsync -av --progress dirac:/Volumes/Space/archive/mirror/ted/ted.20*.sql.bz2 ./data/archive/mirror/ted/
rsync -av --progress dirac:/Volumes/Space/archive/mirror/ted/ted.watt*.sql.bz2 ./data/archive/mirror/ted/
```

### list of restores
```bash
# 2020-12-11T19:43:35.600Z Phase 1- verify start
51	/Volumes/Space/archive/mirror/ted/ted.20090214.1756.sql.bz2
51	/Volumes/Space/archive/mirror/ted/ted.20090214.1800.sql.bz2
64	/Volumes/Space/archive/mirror/ted/ted.20090326.1052.sql.bz2
65	/Volumes/Space/archive/mirror/ted/ted.20090328.1335.sql.bz2
100	/Volumes/Space/archive/mirror/ted/ted.20090528.0815.sql.bz2
100	/Volumes/Space/archive/mirror/ted/ted.20090609.0858.sql.bz2
132	/Volumes/Space/archive/mirror/ted/ted.20090918.0240.sql.bz2
# 2020-12-11T23:06:45.155Z start
148	/Volumes/Space/archive/mirror/ted/ted.20091022.0301.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.20091102.0131.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.20091113.2035.sql.bz2
# 2020-12-12T01:26:06.140Z start
308	/Volumes/Space/archive/mirror/ted/ted.20110406.0317.sql.bz2
324	/Volumes/Space/archive/mirror/ted/ted.20110607.0118.sql.bz2
372	/Volumes/Space/archive/mirror/ted/ted.20111017.2034.sql.bz2
452	/Volumes/Space/archive/mirror/ted/ted.20120608.0122.sql.bz2
532	/Volumes/Space/archive/mirror/ted/ted.20130221.2122.sql.bz2
660	/Volumes/Space/archive/mirror/ted/ted.20140219.2021.sql.bz2
708	/Volumes/Space/archive/mirror/ted/ted.20140806.0019.sql.bz2
836	/Volumes/Space/archive/mirror/ted/ted.20150928.1006.sql.bz2
2020-12-12T20:23:26.248Z done - phase 1 took 24hours

# 2020-12-13T09:35:54.080Z Phase 2 - verify start
7   /Volumes/Space/archive/mirror/ted/ted.watt-just2016.2016-02-14-1624.sql.gz
# exception .gz -> .bz2: ted.watt-just2016.2016-02-14-1624.sql.gz

# 2020-12-13T09:42:21.909Z done first restore
63	/Volumes/Space/archive/mirror/ted/ted.watt.20090918.0300.sql.bz2
68	/Volumes/Space/archive/mirror/ted/ted.watt.20091022.0258.sql.bz2
84	/Volumes/Space/archive/mirror/ted/ted.watt.20091102.0134.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.watt.20110406.0316.sql.bz2
164	/Volumes/Space/archive/mirror/ted/ted.watt.20110607.0115.sql.bz2
180	/Volumes/Space/archive/mirror/ted/ted.watt.20111017.2045.sql.bz2
212	/Volumes/Space/archive/mirror/ted/ted.watt.20120608.0119.sql.bz2
# 2020-12-13T15:06:14.638Z start
260	/Volumes/Space/archive/mirror/ted/ted.watt.20130221.2119.sql.bz2
308	/Volumes/Space/archive/mirror/ted/ted.watt.20140219.2038.sql.bz2
340	/Volumes/Space/archive/mirror/ted/ted.watt.20140806.0016.sql.bz2
340	/Volumes/Space/archive/mirror/ted/ted.watt.20141005.2218.sql.bz2
# 2020-12-13T19:52:16.739Z start
388	/Volumes/Space/archive/mirror/ted/ted.watt.20150928.1003.sql.bz2
420	/Volumes/Space/archive/mirror/ted/ted.watt.2016-02-14-1555.sql.bz2

8   /Volumes/Space/archive/mirror/ted/ted.watt.20160430.0232Z.sql.bz2
16	/Volumes/Space/archive/mirror/ted/ted.watt.20160616.0229Z.sql.bz2
21	/Volumes/Space/archive/mirror/ted/ted.watt.20160719.1848Z.sql.bz2
31	/Volumes/Space/archive/mirror/ted/ted.watt.20160918.0059Z.sql.bz2
43	/Volumes/Space/archive/mirror/ted/ted.watt.20161202.0733Z.sql.bz2
49	/Volumes/Space/archive/mirror/ted/ted.watt.20170106.0629Z.sql.bz2
63	/Volumes/Space/archive/mirror/ted/ted.watt.20170326.1528Z.sql.bz2
84	/Volumes/Space/archive/mirror/ted/ted.watt.20170607.0541Z.sql.bz2
84	/Volumes/Space/archive/mirror/ted/ted.watt.20170727.1724Z.sql.bz2
116	/Volumes/Space/archive/mirror/ted/ted.watt.20180217.2219Z.sql.bz2
132	/Volumes/Space/archive/mirror/ted/ted.watt.20180326.0312Z.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.watt.20180612.0035Z.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.watt.20180720.2138Z.sql.bz2
148	/Volumes/Space/archive/mirror/ted/ted.watt.20180831.2033Z.sql.bz2
164	/Volumes/Space/archive/mirror/ted/ted.watt.20181024.1913Z.sql.bz2
# First time entries for 2007-08-28 appear ^ These were removed in the final rollup
195	/Volumes/Space/archive/mirror/ted/ted.watt.20190414.0128Z.sql.bz2
201	/Volumes/Space/archive/mirror/ted/ted.watt.20190617.0443Z.sql.bz2
209	/Volumes/Space/archive/mirror/ted/ted.watt.20190818.0554Z.sql.bz2
233	/Volumes/Space/archive/mirror/ted/ted.watt.20191129.0710Z.sql.bz2
# 2020-12-14T05:57:44.669Z start
251	/Volumes/Space/archive/mirror/ted/ted.watt.20200413.1503Z.sql.bz2
269	/Volumes/Space/archive/mirror/ted/ted.watt.20200807.2218Z.sql.bz2
# Failed, had to restart
297	/Volumes/Space/archive/mirror/ted/ted.watt.20201120.2332Z.sql.bz2
# Phase 2 verification took 22 hours
```

### /archive/mirror/ted/ted.20090214.1756.sql.bz2
```
- Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090214.1756.sql.bz2
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-02-14 22:55:47	15495136
- Done Restoring Database
real	9m6.224s

JSON Size: 660M
Summary of mysql restore: 

2020-12-10T22:01:13.576Z - Verified jsonl <-> mysql(watt):
2020-12-10T22:01:13.576Z - [2008-07-30T00:04:40Z, 2009-02-14T22:55:47Z](15495136) Equal
2020-12-10T22:01:13.576Z - -=- jsonl <-> mysql(ted_native)
2020-12-10T22:01:56.657Z - jsonl <-> mysql(ted_native) took 43.081s, rate ~ 114.5k/s count: 4934047
2020-12-10T22:01:56.657Z - Verified jsonl <-> mysql(ted_native):
2020-12-10T22:01:56.657Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-10T22:01:56.657Z - [2008-12-17T19:37:16Z, 2009-02-14T22:55:47Z](4934047) Equal
```

#### analysis

```
mysql> select min(stamp),max(stamp),count(*) from ted_native;
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2008-12-17 19:37:16 | 2009-02-14 22:55:47 |  4934047 |
+---------------------+---------------------+----------+
1 row in set (0.00 sec)

mysql> select min(stamp),max(stamp),count(*) from ted_service;
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2008-11-14 23:18:13 | 2008-12-17 19:19:20 |  2710592 |
+---------------------+---------------------+----------+
1 row in set (0.00 sec)


mysql> select min(stamp),max(stamp),count(*) from watt;
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2008-07-30 00:04:40 | 2009-02-14 22:55:47 | 15495136 |
+---------------------+---------------------+----------+

4934047+2710592 = 7644639

mysql> select min(stamp),max(stamp),count(*) from watt where stamp>'2008-12-17 19:19:20';
+---------------------+---------------------+----------+
| min(stamp)          | max(stamp)          | count(*) |
+---------------------+---------------------+----------+
| 2008-12-17 19:37:16 | 2009-02-14 22:55:47 |  4934047 |
+---------------------+---------------------+----------+
4934047+2710592 = 7644639 + 4934047
```

### /archive/mirror/ted/ted.20150928.1006.sql.bz2

- Restoring database from snapshot: ./data/archive/mirror/ted/ted.20150928.1006.sql.bz2
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2015-09-28 14:06:52	212737945
- Done Restoring Database
real	162m43.936s



### ted.watt.2016-02-14-1555.sql.bz2

- Restoring database from snapshot: /Users/daniel/Downloads/ted/ted.watt.2016-02-14-1555.sql.bz2
  To docker container named: go-ted1k_mysql_1
  Using database: ted, MYSQL_USER=ted
  Data Volume will persisted inside that docker container

- Verifying docker environment
  Docker seems to be setup properly

- Waiting for database server (go-ted1k_mysql_1) to accept connections (max 30 seconds)
  Connected
- Restoring database...

- Expect something recent in watt table
  min(stamp) max(stamp) count(\*)
  2008-07-30 00:04:40 2016-02-11 12:22:45 223101124

- Done Restoring Database

  2976.253s

In `ted.watt.2016-02-14-1555.sql.bz2`, there are 10,561,089 samples before the TEDNative cutoff,
whereas the last sqlite snapshot `/archive/mirror/ted/footprints/end-of-life/TED.db.20081217.1419.bz2` has only 10,242,588 samples, so it seems there are 318,501 more samples in the mysql snapshot over this period. Hmm.

mysql> select min(stamp),max(stamp),count(_) from watt where stamp<'2008-12-17 19:19:30';
+---------------------+---------------------+----------+
| min(stamp) | max(stamp) | count(_) |
+---------------------+---------------------+----------+
| 2008-07-30 00:04:40 | 2008-12-17 19:19:20 | 10561089 |
+---------------------+---------------------+----------+

Once restored into Mysql:8.0 this restore took ~6.3Gb

| size | path               |
| ---- | ------------------ |
| 8811 | /var/lib/mysql/    |
| 6385 | /var/lib/mysql/ted |

### exception .gz: ted.watt-just2016.2016-02-14-1624.sql.gz

## Oldest TedNative sqlite dump

- Data capture started at 2008-07-30T00:04:40Z
- Switchover from TedNative (sqlite3) to python v1 happened at 2008-12-17T19:37:16Z
- Note: there is actually more data in mysqldump than in TED's Footprint SQLite3 db, we could inspect the other monthly's in `/archive/mirror/ted/footprints/monnthly/TED.db.2008XX01.0000.bz2`

| stamp               | watt |                                 |
| ------------------- | ---- | ------------------------------- |
| 2008-07-30 00:04:40 | 540  |
| ....                | ...  |
| 2008-12-17 19:19:19 | 570  |
| 2008-12-17 19:19:20 | 560  | <- last sqlite TEDNative insert |
| 2008-12-17 19:37:16 | 610  | <- first mysql python insert    |
| 2008-12-17 19:37:17 | 610  |

From `TED.db.20081217.1419.sqlite` <- `/archive/mirror/ted/footprints/end-of-life/TED.db.20081217.1419.bz2`,
and from `.../im-ted1k/legacy/scalr-utils/scalr.py`,
and finally `docker run --rm -it python:2.7`

```sqlite
sqlite> select min(tick),max(tick),count(*) from rdu_second_data;
0633529550800006250|0633651203640000000|10242588
```

```python
import string
def tedToSecs(tedTimeString):
  millis = string.atol(tedTimeString)/10000
  return millis / 1000 - 62135578800;

import time
def tedToLocal(tedTimeString):
  secs = tedToSecs(tedTimeString)
  return time.strftime("%Y-%m-%d %H:%M:%S %Z",time.localtime(secs))

print tedToSecs('0633529550800006250') # 1217376280
print tedToSecs('0633651203640000000')# 1229541564

print tedToLocal('0633529550800006250') # 2008-07-30 00:04:40 UTC
print tedToLocal('0633651203640000000')# 2008-12-17 19:19:24 UTC
```

## Restore phase 1 - verification - detailed logs

```go
verify("jsonl <-> mysql(watt)", jsonl.NewReader(), mysql.NewReader(db, "watt"))
verify("jsonl <-> mysql(ted_native)", jsonl.NewReader(), mysql.NewReader(db, "ted_native"))
```

```bash
~/go-ted1k$ time ./restore-db.sh 
- Verifying docker environment
WARNING: No swap limit support
  Docker seems to be setup properly

- Waiting for database server (go-ted1k_mysql_1) to accept connections (max 30 seconds)
  Connected

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090214.1756.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	9m7.466s
user	0m16.462s
sys	0m4.937s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-02-14 22:55:47	4934047
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-02-14 22:55:47	15495136

- run mysql restore
2020-12-11T19:53:35.600Z - Starting TED1K mysql restore
2020-12-11T19:53:35.602Z - Connected to MySQL
2020-12-11T19:53:35.602Z - -=- jsonl <-> mysql(watt)
2020-12-11T19:54:58.587Z - jsonl <-> mysql(watt) took 1m22.985s, rate ~ 186.7k/s count: 15495136
2020-12-11T20:00:48.723Z - Verified jsonl <-> mysql(watt):
2020-12-11T20:00:48.723Z - [2008-07-30T00:04:40Z, 2009-02-14T22:55:47Z](15495136) Equal
2020-12-11T20:00:48.723Z - [2009-02-14T22:55:48Z, 2015-09-28T14:06:52Z](197242834) MissingInB
2020-12-11T20:00:48.723Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T20:01:33.334Z - jsonl <-> mysql(ted_native) took 44.611s, rate ~ 110.6k/s count: 4934047
2020-12-11T20:07:27.297Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T20:07:27.297Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T20:07:27.297Z - [2008-12-17T19:37:16Z, 2009-02-14T22:55:47Z](4934047) Equal
2020-12-11T20:07:27.297Z - [2009-02-14T22:55:48Z, 2015-09-28T14:06:52Z](197242834) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090214.1800.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	9m3.722s
user	0m16.089s
sys	0m4.846s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-02-14 23:00:05	4934304
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-02-14 23:00:05	15495393

- run mysql restore
2020-12-11T20:16:37.288Z - Starting TED1K mysql restore
2020-12-11T20:16:37.309Z - Connected to MySQL
2020-12-11T20:16:37.309Z - -=- jsonl <-> mysql(watt)
2020-12-11T20:17:58.138Z - jsonl <-> mysql(watt) took 1m20.829s, rate ~ 191.7k/s count: 15495393
2020-12-11T20:23:47.680Z - Verified jsonl <-> mysql(watt):
2020-12-11T20:23:47.680Z - [2008-07-30T00:04:40Z, 2009-02-14T23:00:05Z](15495393) Equal
2020-12-11T20:23:47.680Z - [2009-02-14T23:00:06Z, 2015-09-28T14:06:52Z](197242577) MissingInB
2020-12-11T20:23:47.680Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T20:24:32.667Z - jsonl <-> mysql(ted_native) took 44.987s, rate ~ 109.7k/s count: 4934304
2020-12-11T20:30:24.277Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T20:30:24.277Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T20:30:24.277Z - [2008-12-17T19:37:16Z, 2009-02-14T23:00:05Z](4934304) Equal
2020-12-11T20:30:24.277Z - [2009-02-14T23:00:06Z, 2015-09-28T14:06:52Z](197242577) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090326.1052.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	11m54.704s
user	0m21.218s
sys	0m6.202s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-03-26 14:52:59	8314404
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-03-26 14:52:59	18875493

- run mysql restore
2020-12-11T20:42:24.360Z - Starting TED1K mysql restore
2020-12-11T20:42:24.383Z - Connected to MySQL
2020-12-11T20:42:24.383Z - -=- jsonl <-> mysql(watt)
2020-12-11T20:44:05.108Z - jsonl <-> mysql(watt) took 1m40.725s, rate ~ 187.4k/s count: 18875493
2020-12-11T20:49:53.188Z - Verified jsonl <-> mysql(watt):
2020-12-11T20:49:53.188Z - [2008-07-30T00:04:40Z, 2009-03-26T14:52:59Z](18875493) Equal
2020-12-11T20:49:53.188Z - [2009-03-26T14:53:00Z, 2015-09-28T14:06:52Z](193862477) MissingInB
2020-12-11T20:49:53.188Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T20:50:56.180Z - jsonl <-> mysql(ted_native) took 1m2.992s, rate ~ 132.0k/s count: 8314404
2020-12-11T20:56:43.698Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T20:56:43.698Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T20:56:43.698Z - [2008-12-17T19:37:16Z, 2009-03-26T14:52:59Z](8314404) Equal
2020-12-11T20:56:43.698Z - [2009-03-26T14:53:00Z, 2015-09-28T14:06:52Z](193862477) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090328.1335.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	12m11.046s
user	0m21.024s
sys	0m6.546s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-03-28 17:35:27	8494287
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-03-28 17:35:27	19055376

- run mysql restore
2020-12-11T21:08:59.818Z - Starting TED1K mysql restore
2020-12-11T21:08:59.844Z - Connected to MySQL
2020-12-11T21:08:59.844Z - -=- jsonl <-> mysql(watt)
2020-12-11T21:10:42.171Z - jsonl <-> mysql(watt) took 1m42.326s, rate ~ 186.2k/s count: 19055376
2020-12-11T21:16:31.380Z - Verified jsonl <-> mysql(watt):
2020-12-11T21:16:31.380Z - [2008-07-30T00:04:40Z, 2009-03-28T17:35:27Z](19055376) Equal
2020-12-11T21:16:31.380Z - [2009-03-28T17:35:28Z, 2015-09-28T14:06:52Z](193682594) MissingInB
2020-12-11T21:16:31.380Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T21:17:35.693Z - jsonl <-> mysql(ted_native) took 1m4.313s, rate ~ 132.1k/s count: 8494287
2020-12-11T21:23:24.028Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T21:23:24.028Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T21:23:24.028Z - [2008-12-17T19:37:16Z, 2009-03-28T17:35:27Z](8494287) Equal
2020-12-11T21:23:24.028Z - [2009-03-28T17:35:28Z, 2015-09-28T14:06:52Z](193682594) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090528.0815.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	16m11.319s
user	0m28.646s
sys	0m8.181s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-05-28 12:15:46	13668586
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-05-28 12:15:46	24229675

- run mysql restore
2020-12-11T21:39:40.675Z - Starting TED1K mysql restore
2020-12-11T21:39:40.735Z - Connected to MySQL
2020-12-11T21:39:40.735Z - -=- jsonl <-> mysql(watt)
2020-12-11T21:41:51.377Z - jsonl <-> mysql(watt) took 2m10.641s, rate ~ 185.5k/s count: 24229675
2020-12-11T21:47:27.620Z - Verified jsonl <-> mysql(watt):
2020-12-11T21:47:27.620Z - [2008-07-30T00:04:40Z, 2009-05-28T12:15:46Z](24229675) Equal
2020-12-11T21:47:27.620Z - [2009-05-28T12:15:47Z, 2015-09-28T14:06:52Z](188508295) MissingInB
2020-12-11T21:47:27.620Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T21:49:00.478Z - jsonl <-> mysql(ted_native) took 1m32.858s, rate ~ 147.2k/s count: 13668586
2020-12-11T21:54:36.027Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T21:54:36.027Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T21:54:36.028Z - [2008-12-17T19:37:16Z, 2009-05-28T12:15:46Z](13668586) Equal
2020-12-11T21:54:36.028Z - [2009-05-28T12:15:47Z, 2015-09-28T14:06:52Z](188508295) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090609.0858.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	16m52.486s
user	0m30.672s
sys	0m8.975s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-06-09 12:58:55	14697489
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-06-09 12:58:55	25258578

- run mysql restore
2020-12-11T22:11:33.552Z - Starting TED1K mysql restore
2020-12-11T22:11:33.578Z - Connected to MySQL
2020-12-11T22:11:33.578Z - -=- jsonl <-> mysql(watt)
2020-12-11T22:13:48.434Z - jsonl <-> mysql(watt) took 2m14.855s, rate ~ 187.3k/s count: 25258578
2020-12-11T22:19:24.506Z - Verified jsonl <-> mysql(watt):
2020-12-11T22:19:24.506Z - [2008-07-30T00:04:40Z, 2009-06-09T12:58:55Z](25258578) Equal
2020-12-11T22:19:24.506Z - [2009-06-09T12:58:56Z, 2015-09-28T14:06:52Z](187479392) MissingInB
2020-12-11T22:19:24.506Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T22:21:01.608Z - jsonl <-> mysql(ted_native) took 1m37.102s, rate ~ 151.4k/s count: 14697489
2020-12-11T22:26:36.025Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T22:26:36.026Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T22:26:36.026Z - [2008-12-17T19:37:16Z, 2009-06-09T12:58:55Z](14697489) Equal
2020-12-11T22:26:36.026Z - [2009-06-09T12:58:56Z, 2015-09-28T14:06:52Z](187479392) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20090918.0240.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	23m59.822s
user	0m40.754s
sys	0m12.476s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-09-18 06:40:17	23312671
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-09-18 06:40:17	33873760

- run mysql restore
2020-12-11T22:50:40.882Z - Starting TED1K mysql restore
2020-12-11T22:50:40.912Z - Connected to MySQL
2020-12-11T22:50:40.912Z - -=- jsonl <-> mysql(watt)
2020-12-11T22:53:31.467Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 184.9k/s, took 2m50.554s, rate ~ 184.9k/s count: 31536000
2020-12-11T22:53:44.116Z - jsonl <-> mysql(watt) took 3m3.203s, rate ~ 184.9k/s count: 33873760
2020-12-11T22:59:04.371Z - Verified jsonl <-> mysql(watt):
2020-12-11T22:59:04.371Z - [2008-07-30T00:04:40Z, 2009-09-18T06:40:17Z](33873760) Equal
2020-12-11T22:59:04.371Z - [2009-09-18T06:40:18Z, 2015-09-28T14:06:52Z](178864210) MissingInB
2020-12-11T22:59:04.371Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T23:01:26.657Z - jsonl <-> mysql(ted_native) took 2m22.286s, rate ~ 163.8k/s count: 23312671
2020-12-11T23:06:45.155Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T23:06:45.155Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T23:06:45.155Z - [2008-12-17T19:37:16Z, 2009-09-18T06:40:17Z](23312671) Equal
2020-12-11T23:06:45.155Z - [2009-09-18T06:40:18Z, 2015-09-28T14:06:52Z](178864210) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20091022.0301.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	27m1.348s
user	0m46.351s
sys	0m14.627s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-10-22 07:01:15	26221838
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-10-22 07:01:15	36782927

- run mysql restore
2020-12-11T23:33:51.864Z - Starting TED1K mysql restore
2020-12-11T23:33:51.874Z - Connected to MySQL
2020-12-11T23:33:51.874Z - -=- jsonl <-> mysql(watt)
2020-12-11T23:37:16.448Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 154.2k/s, took 3m24.573s, rate ~ 154.2k/s count: 31536000
2020-12-11T23:37:50.443Z - jsonl <-> mysql(watt) took 3m58.568s, rate ~ 154.2k/s count: 36782927
^[f2020-12-11T23:43:06.840Z - Verified jsonl <-> mysql(watt):
2020-12-11T23:43:06.840Z - [2008-07-30T00:04:40Z, 2009-10-22T07:01:15Z](36782927) Equal
2020-12-11T23:43:06.840Z - [2009-10-22T20:13:02Z, 2015-09-28T14:06:52Z](175955043) MissingInB
2020-12-11T23:43:06.840Z - -=- jsonl <-> mysql(ted_native)
2020-12-11T23:46:10.781Z - jsonl <-> mysql(ted_native) took 3m3.94s, rate ~ 142.6k/s count: 26221838
2020-12-11T23:51:28.072Z - Verified jsonl <-> mysql(ted_native):
2020-12-11T23:51:28.072Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-11T23:51:28.072Z - [2008-12-17T19:37:16Z, 2009-10-22T07:01:15Z](26221838) Equal
2020-12-11T23:51:28.072Z - [2009-10-22T20:13:02Z, 2015-09-28T14:06:52Z](175955043) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20091102.0131.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	29m0.210s
user	0m46.468s
sys	0m15.921s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-11-02 06:31:14	27064883
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-11-02 06:31:14	37625972

- run mysql restore
2020-12-12T00:20:33.511Z - Starting TED1K mysql restore
2020-12-12T00:20:33.513Z - Connected to MySQL
2020-12-12T00:20:33.513Z - -=- jsonl <-> mysql(watt)
2020-12-12T00:23:56.071Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.7k/s, took 3m22.557s, rate ~ 155.7k/s count: 31536000
2020-12-12T00:24:34.853Z - jsonl <-> mysql(watt) took 4m1.34s, rate ~ 155.9k/s count: 37625972
2020-12-12T00:29:50.750Z - Verified jsonl <-> mysql(watt):
2020-12-12T00:29:50.750Z - [2008-07-30T00:04:40Z, 2009-11-02T06:31:14Z](37625972) Equal
2020-12-12T00:29:50.750Z - [2009-11-02T06:31:15Z, 2015-09-28T14:06:52Z](175111998) MissingInB
2020-12-12T00:29:50.750Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T00:33:01.838Z - jsonl <-> mysql(ted_native) took 3m11.088s, rate ~ 141.6k/s count: 27064883
2020-12-12T00:38:17.164Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T00:38:17.164Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T00:38:17.164Z - [2008-12-17T19:37:16Z, 2009-11-02T06:31:14Z](27064883) Equal
2020-12-12T00:38:17.164Z - [2009-11-02T06:31:15Z, 2015-09-28T14:06:52Z](175111998) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20091113.2035.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	29m50.511s
user	0m47.759s
sys	0m16.531s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2009-11-14 01:35:35	28066796
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-11-14 01:35:35	38627885

- run mysql restore
2020-12-12T01:08:13.039Z - Starting TED1K mysql restore
2020-12-12T01:08:13.041Z - Connected to MySQL
2020-12-12T01:08:13.041Z - -=- jsonl <-> mysql(watt)
2020-12-12T01:11:35.395Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.8k/s, took 3m22.355s, rate ~ 155.8k/s count: 31536000
2020-12-12T01:12:20.750Z - jsonl <-> mysql(watt) took 4m7.709s, rate ~ 155.9k/s count: 38627885
2020-12-12T01:17:35.557Z - Verified jsonl <-> mysql(watt):
2020-12-12T01:17:35.557Z - [2008-07-30T00:04:40Z, 2009-11-14T01:35:35Z](38627885) Equal
2020-12-12T01:17:35.557Z - [2009-11-14T01:48:36Z, 2015-09-28T14:06:52Z](174110085) MissingInB
2020-12-12T01:17:35.557Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T01:20:51.538Z - jsonl <-> mysql(ted_native) took 3m15.98s, rate ~ 143.2k/s count: 28066796
2020-12-12T01:26:06.140Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T01:26:06.140Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T01:26:06.140Z - [2008-12-17T19:37:16Z, 2009-11-14T01:35:35Z](28066796) Equal
2020-12-12T01:26:06.140Z - [2009-11-14T01:48:36Z, 2015-09-28T14:06:52Z](174110085) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20110406.0317.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	63m33.269s
user	1m41.827s
sys	0m35.158s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2011-04-06 07:15:35	67085417
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2011-04-06 07:15:35	77646482

- run mysql restore
2020-12-12T02:29:44.824Z - Starting TED1K mysql restore
2020-12-12T02:29:44.827Z - Connected to MySQL
2020-12-12T02:29:44.827Z - -=- jsonl <-> mysql(watt)
2020-12-12T02:33:07.631Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.5k/s, took 3m22.804s, rate ~ 155.5k/s count: 31536000
2020-12-12T02:36:28.536Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 157.0k/s, took 6m43.709s, rate ~ 156.2k/s count: 63072000
2020-12-12T02:38:01.010Z - jsonl <-> mysql(watt) took 8m16.183s, rate ~ 156.5k/s count: 77646482
2020-12-12T02:42:02.652Z - Verified jsonl <-> mysql(watt):
2020-12-12T02:42:02.652Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T02:42:02.652Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T02:42:02.652Z - [2009-11-17T03:00:30Z, 2011-04-06T07:15:35Z](38760401) Equal
2020-12-12T02:42:02.652Z - [2011-04-06T07:24:52Z, 2015-09-28T14:06:52Z](135091463) MissingInB
2020-12-12T02:42:02.652Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T02:45:42.851Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 143.2k/s, took 3m40.199s, rate ~ 143.2k/s count: 31536000
2020-12-12T02:49:02.894Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.6k/s, took 7m0.241s, rate ~ 150.1k/s count: 63072000
2020-12-12T02:49:28.212Z - jsonl <-> mysql(ted_native) took 7m25.56s, rate ~ 150.6k/s count: 67085417
2020-12-12T02:53:30.145Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T02:53:30.145Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T02:53:30.145Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T02:53:30.145Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T02:53:30.145Z - [2009-11-17T03:02:49Z, 2011-04-06T07:15:35Z](38760400) Equal
2020-12-12T02:53:30.145Z - [2011-04-06T07:24:52Z, 2015-09-28T14:06:52Z](135091463) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20110607.0118.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	67m16.029s
user	1m47.745s
sys	0m36.741s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2011-06-07 05:18:07	72370533
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2011-06-07 05:18:07	82931598

- run mysql restore
2020-12-12T04:00:51.521Z - Starting TED1K mysql restore
2020-12-12T04:00:51.523Z - Connected to MySQL
2020-12-12T04:00:51.523Z - -=- jsonl <-> mysql(watt)
2020-12-12T04:04:13.991Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.8k/s, took 3m22.467s, rate ~ 155.8k/s count: 31536000
2020-12-12T04:07:35.114Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 156.8k/s, took 6m43.59s, rate ~ 156.3k/s count: 63072000
2020-12-12T04:09:41.351Z - jsonl <-> mysql(watt) took 8m49.827s, rate ~ 156.5k/s count: 82931598
2020-12-12T04:13:35.374Z - Verified jsonl <-> mysql(watt):
2020-12-12T04:13:35.374Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T04:13:35.374Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T04:13:35.374Z - [2009-11-17T03:00:30Z, 2011-06-07T05:18:07Z](44045517) Equal
2020-12-12T04:13:35.374Z - [2011-06-07T05:18:08Z, 2015-09-28T14:06:52Z](129806347) MissingInB
2020-12-12T04:13:35.374Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T04:17:15.020Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 143.6k/s, took 3m39.645s, rate ~ 143.6k/s count: 31536000
2020-12-12T04:20:36.332Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 156.7k/s, took 7m0.957s, rate ~ 149.8k/s count: 63072000
2020-12-12T04:21:35.184Z - jsonl <-> mysql(ted_native) took 7m59.809s, rate ~ 150.8k/s count: 72370533
2020-12-12T04:25:28.341Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T04:25:28.341Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T04:25:28.341Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T04:25:28.341Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T04:25:28.341Z - [2009-11-17T03:02:49Z, 2011-06-07T05:18:07Z](44045516) Equal
2020-12-12T04:25:28.341Z - [2011-06-07T05:18:08Z, 2015-09-28T14:06:52Z](129806347) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20111017.2034.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	75m3.564s
user	2m0.566s
sys	0m41.336s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2011-10-18 00:33:13	81986218
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2011-10-18 00:33:13	92547284

- run mysql restore
2020-12-12T05:40:37.456Z - Starting TED1K mysql restore
2020-12-12T05:40:37.458Z - Connected to MySQL
2020-12-12T05:40:37.458Z - -=- jsonl <-> mysql(watt)
2020-12-12T05:44:01.130Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 154.8k/s, took 3m23.671s, rate ~ 154.8k/s count: 31536000
2020-12-12T05:47:24.944Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 154.7k/s, took 6m47.485s, rate ~ 154.8k/s count: 63072000
2020-12-12T05:50:35.005Z - jsonl <-> mysql(watt) took 9m57.547s, rate ~ 154.9k/s count: 92547284
2020-12-12T05:54:10.505Z - Verified jsonl <-> mysql(watt):
2020-12-12T05:54:10.505Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T05:54:10.505Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T05:54:10.505Z - [2009-11-17T03:00:30Z, 2011-10-18T00:33:13Z](53661203) Equal
2020-12-12T05:54:10.505Z - [2011-10-18T00:51:14Z, 2015-09-28T14:06:52Z](120190661) MissingInB
2020-12-12T05:54:10.505Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T05:57:50.793Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 143.2k/s, took 3m40.287s, rate ~ 143.2k/s count: 31536000
2020-12-12T06:01:11.576Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.1k/s, took 7m1.071s, rate ~ 149.8k/s count: 63072000
2020-12-12T06:03:10.927Z - jsonl <-> mysql(ted_native) took 9m0.422s, rate ~ 151.7k/s count: 81986218
2020-12-12T06:06:45.631Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T06:06:45.631Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T06:06:45.631Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T06:06:45.631Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T06:06:45.631Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T06:06:45.631Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T06:06:45.631Z - [2011-10-16T21:54:25Z, 2011-10-18T00:33:13Z](94912) Equal
2020-12-12T06:06:45.632Z - [2011-10-18T00:51:14Z, 2015-09-28T14:06:52Z](120190661) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20120608.0122.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	93m42.371s
user	2m28.268s
sys	0m51.801s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2012-06-08 05:22:10	101938261
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2012-06-08 05:22:10	112499327

- run mysql restore
2020-12-12T07:40:33.556Z - Starting TED1K mysql restore
2020-12-12T07:40:33.558Z - Connected to MySQL
2020-12-12T07:40:33.558Z - -=- jsonl <-> mysql(watt)
2020-12-12T07:43:52.555Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 158.5k/s, took 3m18.997s, rate ~ 158.5k/s count: 31536000
2020-12-12T07:47:09.874Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 159.8k/s, took 6m36.316s, rate ~ 159.1k/s count: 63072000
2020-12-12T07:50:26.057Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 160.7k/s, took 9m52.499s, rate ~ 159.7k/s count: 94608000
2020-12-12T07:52:16.787Z - jsonl <-> mysql(watt) took 11m43.229s, rate ~ 160.0k/s count: 112499327
2020-12-12T07:55:17.159Z - Verified jsonl <-> mysql(watt):
2020-12-12T07:55:17.159Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T07:55:17.159Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T07:55:17.159Z - [2009-11-17T03:00:30Z, 2012-06-08T05:22:10Z](73613246) Equal
2020-12-12T07:55:17.159Z - [2012-06-08T05:22:12Z, 2015-09-28T14:06:52Z](100238618) MissingInB
2020-12-12T07:55:17.159Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T07:58:55.970Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 144.1k/s, took 3m38.811s, rate ~ 144.1k/s count: 31536000
2020-12-12T08:02:16.134Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.6k/s, took 6m58.975s, rate ~ 150.5k/s count: 63072000
2020-12-12T08:05:35.446Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 158.2k/s, took 10m18.287s, rate ~ 153.0k/s count: 94608000
2020-12-12T08:06:21.455Z - jsonl <-> mysql(ted_native) took 11m4.296s, rate ~ 153.5k/s count: 101938261
2020-12-12T08:09:20.793Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T08:09:20.793Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T08:09:20.793Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T08:09:20.793Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T08:09:20.793Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T08:09:20.793Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T08:09:20.793Z - [2011-10-16T21:54:25Z, 2012-06-08T05:22:10Z](20046955) Equal
2020-12-12T08:09:20.793Z - [2012-06-08T05:22:12Z, 2015-09-28T14:06:52Z](100238618) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20130221.2122.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	112m54.141s
user	2m57.062s
sys	1m1.232s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2013-02-22 02:22:26	123979218
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2013-02-22 02:22:26	134540284

- run mysql restore
2020-12-12T10:02:20.797Z - Starting TED1K mysql restore
2020-12-12T10:02:20.799Z - Connected to MySQL
2020-12-12T10:02:20.799Z - -=- jsonl <-> mysql(watt)
2020-12-12T10:05:42.034Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 156.7k/s, took 3m21.235s, rate ~ 156.7k/s count: 31536000
2020-12-12T10:09:02.179Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 157.6k/s, took 6m41.38s, rate ~ 157.1k/s count: 63072000
2020-12-12T10:12:21.672Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 158.1k/s, took 10m0.873s, rate ~ 157.5k/s count: 94608000
2020-12-12T10:15:40.255Z - jsonl <-> mysql(watt) (2012-11-15) inner rate ~ 158.8k/s, took 13m19.456s, rate ~ 157.8k/s count: 126144000
2020-12-12T10:16:33.353Z - jsonl <-> mysql(watt) took 14m12.555s, rate ~ 157.8k/s count: 134540284
2020-12-12T10:18:53.624Z - Verified jsonl <-> mysql(watt):
2020-12-12T10:18:53.624Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T10:18:53.624Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T10:18:53.624Z - [2009-11-17T03:00:30Z, 2013-02-22T02:22:26Z](95654203) Equal
2020-12-12T10:18:53.624Z - [2013-02-22T02:22:27Z, 2015-09-28T14:06:52Z](78197661) MissingInB
2020-12-12T10:18:53.624Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T10:22:32.904Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 143.8k/s, took 3m39.28s, rate ~ 143.8k/s count: 31536000
2020-12-12T10:25:52.975Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.6k/s, took 6m59.351s, rate ~ 150.4k/s count: 63072000
2020-12-12T10:29:12.390Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 158.1k/s, took 10m18.766s, rate ~ 152.9k/s count: 94608000
2020-12-12T10:32:18.052Z - jsonl <-> mysql(ted_native) took 13m24.428s, rate ~ 154.1k/s count: 123979218
2020-12-12T10:34:38.073Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T10:34:38.073Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T10:34:38.073Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T10:34:38.073Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T10:34:38.073Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T10:34:38.073Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T10:34:38.073Z - [2011-10-16T21:54:25Z, 2013-02-22T02:22:26Z](42087912) Equal
2020-12-12T10:34:38.073Z - [2013-02-22T02:22:27Z, 2015-09-28T14:06:52Z](78197661) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20140219.2021.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	137m59.671s
user	3m40.321s
sys	1m16.615s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2014-02-19 13:58:18	154748431
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2014-02-19 13:58:18	165309499

- run mysql restore
2020-12-12T12:52:43.302Z - Starting TED1K mysql restore
2020-12-12T12:52:43.305Z - Connected to MySQL
2020-12-12T12:52:43.305Z - -=- jsonl <-> mysql(watt)
2020-12-12T12:56:05.707Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.8k/s, took 3m22.402s, rate ~ 155.8k/s count: 31536000
2020-12-12T12:59:26.945Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 156.7k/s, took 6m43.64s, rate ~ 156.3k/s count: 63072000
2020-12-12T13:02:48.223Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 156.7k/s, took 10m4.918s, rate ~ 156.4k/s count: 94608000
2020-12-12T13:06:08.606Z - jsonl <-> mysql(watt) (2012-11-15) inner rate ~ 157.4k/s, took 13m25.3s, rate ~ 156.6k/s count: 126144000
2020-12-12T13:09:27.310Z - jsonl <-> mysql(watt) (2013-11-22) inner rate ~ 158.7k/s, took 16m44.004s, rate ~ 157.1k/s count: 157680000
2020-12-12T13:10:15.117Z - jsonl <-> mysql(watt) took 17m31.811s, rate ~ 157.2k/s count: 165309499
2020-12-12T13:11:40.424Z - Verified jsonl <-> mysql(watt):
2020-12-12T13:11:40.424Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T13:11:40.424Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T13:11:40.424Z - [2009-11-17T03:00:30Z, 2014-02-19T13:58:18Z](126423418) Equal
2020-12-12T13:11:40.424Z - [2014-02-20T01:44:19Z, 2015-09-28T14:06:52Z](47428446) MissingInB
2020-12-12T13:11:40.424Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T13:15:19.441Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 144.0k/s, took 3m39.016s, rate ~ 144.0k/s count: 31536000
2020-12-12T13:18:39.670Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 157.5k/s, took 6m59.245s, rate ~ 150.4k/s count: 63072000
2020-12-12T13:21:59.463Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 157.8k/s, took 10m19.038s, rate ~ 152.8k/s count: 94608000
2020-12-12T13:25:20.132Z - jsonl <-> mysql(ted_native) (2013-03-19) inner rate ~ 157.2k/s, took 13m39.708s, rate ~ 153.9k/s count: 126144000
2020-12-12T13:28:21.946Z - jsonl <-> mysql(ted_native) took 16m41.521s, rate ~ 154.5k/s count: 154748431
2020-12-12T13:29:46.838Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T13:29:46.838Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T13:29:46.839Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T13:29:46.839Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T13:29:46.839Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T13:29:46.839Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T13:29:46.839Z - [2011-10-16T21:54:25Z, 2013-10-18T13:53:18Z](62272852) Equal
2020-12-12T13:29:46.839Z - [2013-10-18T15:53:27Z, 2013-10-18T15:58:38Z](2) MissingInB
2020-12-12T13:29:46.839Z - [2013-10-18T16:06:13Z, 2014-02-19T13:58:18Z](10584273) Equal
2020-12-12T13:29:46.839Z - [2014-02-20T01:44:19Z, 2015-09-28T14:06:52Z](47428446) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20140806.0019.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	150m48.392s
user	4m3.086s
sys	1m25.752s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2014-08-06 04:19:27	168924060
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2014-08-06 04:19:27	179485128

- run mysql restore
2020-12-12T16:00:40.893Z - Starting TED1K mysql restore
2020-12-12T16:00:40.895Z - Connected to MySQL
2020-12-12T16:00:40.896Z - -=- jsonl <-> mysql(watt)
2020-12-12T16:04:04.369Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 155.0k/s, took 3m23.473s, rate ~ 155.0k/s count: 31536000
2020-12-12T16:07:26.899Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 155.7k/s, took 6m46.003s, rate ~ 155.3k/s count: 63072000
2020-12-12T16:10:48.511Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 156.4k/s, took 10m7.615s, rate ~ 155.7k/s count: 94608000
2020-12-12T16:14:09.700Z - jsonl <-> mysql(watt) (2012-11-15) inner rate ~ 156.7k/s, took 13m28.804s, rate ~ 156.0k/s count: 126144000
2020-12-12T16:17:31.560Z - jsonl <-> mysql(watt) (2013-11-22) inner rate ~ 156.2k/s, took 16m50.664s, rate ~ 156.0k/s count: 157680000
2020-12-12T16:19:50.132Z - jsonl <-> mysql(watt) took 19m9.237s, rate ~ 156.2k/s count: 179485128
2020-12-12T16:20:50.016Z - Verified jsonl <-> mysql(watt):
2020-12-12T16:20:50.016Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T16:20:50.016Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T16:20:50.016Z - [2009-11-17T03:00:30Z, 2014-08-06T04:19:27Z](140599047) Equal
2020-12-12T16:20:50.016Z - [2014-08-06T04:19:28Z, 2015-09-28T14:06:52Z](33252817) MissingInB
2020-12-12T16:20:50.016Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T16:24:25.925Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 146.1k/s, took 3m35.909s, rate ~ 146.1k/s count: 31536000
2020-12-12T16:27:43.964Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 159.2k/s, took 6m53.948s, rate ~ 152.4k/s count: 63072000
2020-12-12T16:31:01.606Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 159.6k/s, took 10m11.589s, rate ~ 154.7k/s count: 94608000
2020-12-12T16:34:17.981Z - jsonl <-> mysql(ted_native) (2013-03-19) inner rate ~ 160.6k/s, took 13m27.965s, rate ~ 156.1k/s count: 126144000
2020-12-12T16:37:34.224Z - jsonl <-> mysql(ted_native) (2014-03-26) inner rate ~ 160.7k/s, took 16m44.208s, rate ~ 157.0k/s count: 157680000
2020-12-12T16:38:43.575Z - jsonl <-> mysql(ted_native) took 17m53.558s, rate ~ 157.3k/s count: 168924060
2020-12-12T16:39:43.291Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T16:39:43.291Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T16:39:43.291Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T16:39:43.291Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T16:39:43.291Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T16:39:43.291Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T16:39:43.291Z - [2011-10-16T21:54:25Z, 2013-10-18T13:53:18Z](62272852) Equal
2020-12-12T16:39:43.291Z - [2013-10-18T15:53:27Z, 2013-10-18T15:58:38Z](2) MissingInB
2020-12-12T16:39:43.291Z - [2013-10-18T16:06:13Z, 2014-08-06T04:19:27Z](24759902) Equal
2020-12-12T16:39:43.291Z - [2014-08-06T04:19:28Z, 2015-09-28T14:06:52Z](33252817) MissingInB
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.20150928.1006.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...

real	179m23.209s
user	4m44.234s
sys	1m40.917s

- Expect something recent in ted_native table
min(stamp)	max(stamp)	count(*)
2008-12-17 19:37:16	2015-09-28 14:06:52	202176877
- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2015-09-28 14:06:52	212737945

- run mysql restore
2020-12-12T19:39:12.348Z - Starting TED1K mysql restore
2020-12-12T19:39:12.350Z - Connected to MySQL
2020-12-12T19:39:12.350Z - -=- jsonl <-> mysql(watt)
2020-12-12T19:42:32.749Z - jsonl <-> mysql(watt) (2009-08-21) inner rate ~ 157.4k/s, took 3m20.398s, rate ~ 157.4k/s count: 31536000
2020-12-12T19:45:52.706Z - jsonl <-> mysql(watt) (2010-10-16) inner rate ~ 157.7k/s, took 6m40.354s, rate ~ 157.5k/s count: 63072000
2020-12-12T19:49:10.759Z - jsonl <-> mysql(watt) (2011-11-11) inner rate ~ 159.2k/s, took 9m58.407s, rate ~ 158.1k/s count: 94608000
2020-12-12T19:52:29.233Z - jsonl <-> mysql(watt) (2012-11-15) inner rate ~ 158.9k/s, took 13m16.881s, rate ~ 158.3k/s count: 126144000
2020-12-12T19:55:46.680Z - jsonl <-> mysql(watt) (2013-11-22) inner rate ~ 159.7k/s, took 16m34.329s, rate ~ 158.6k/s count: 157680000
2020-12-12T19:59:04.338Z - jsonl <-> mysql(watt) (2014-12-15) inner rate ~ 159.5k/s, took 19m51.987s, rate ~ 158.7k/s count: 189216000
2020-12-12T20:01:32.134Z - jsonl <-> mysql(watt) took 22m19.782s, rate ~ 158.8k/s count: 212737945
2020-12-12T20:01:32.134Z - Verified jsonl <-> mysql(watt):
2020-12-12T20:01:32.134Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-12T20:01:32.134Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInB
2020-12-12T20:01:32.134Z - [2009-11-17T03:00:30Z, 2015-09-28T14:06:52Z](173851864) Equal
2020-12-12T20:01:32.134Z - -=- jsonl <-> mysql(ted_native)
2020-12-12T20:05:13.716Z - jsonl <-> mysql(ted_native) (2009-12-26) inner rate ~ 142.3k/s, took 3m41.582s, rate ~ 142.3k/s count: 31536000
2020-12-12T20:08:35.481Z - jsonl <-> mysql(ted_native) (2011-02-18) inner rate ~ 156.3k/s, took 7m3.347s, rate ~ 149.0k/s count: 63072000
2020-12-12T20:11:57.441Z - jsonl <-> mysql(ted_native) (2012-03-13) inner rate ~ 156.1k/s, took 10m25.307s, rate ~ 151.3k/s count: 94608000
2020-12-12T20:15:19.891Z - jsonl <-> mysql(ted_native) (2013-03-19) inner rate ~ 155.8k/s, took 13m47.757s, rate ~ 152.4k/s count: 126144000
2020-12-12T20:18:41.993Z - jsonl <-> mysql(ted_native) (2014-03-26) inner rate ~ 156.0k/s, took 17m9.859s, rate ~ 153.1k/s count: 157680000
2020-12-12T20:22:03.403Z - jsonl <-> mysql(ted_native) (2015-04-27) inner rate ~ 156.6k/s, took 20m31.269s, rate ~ 153.7k/s count: 189216000
2020-12-12T20:23:26.247Z - jsonl <-> mysql(ted_native) took 21m54.113s, rate ~ 153.9k/s count: 202176877
2020-12-12T20:23:26.248Z - Verified jsonl <-> mysql(ted_native):
2020-12-12T20:23:26.248Z - [2008-07-30T00:04:40Z, 2008-12-17T19:19:20Z](10561089) MissingInB
2020-12-12T20:23:26.248Z - [2008-12-17T19:37:16Z, 2009-11-17T02:46:27Z](28325017) Equal
2020-12-12T20:23:26.248Z - [2009-11-17T03:00:30Z, 2009-11-17T03:00:30Z](1) MissingInB
2020-12-12T20:23:26.248Z - [2009-11-17T03:02:49Z, 2011-10-16T20:56:53Z](53566289) Equal
2020-12-12T20:23:26.248Z - [2011-10-16T21:38:44Z, 2011-10-16T21:38:44Z](1) MissingInB
2020-12-12T20:23:26.248Z - [2011-10-16T21:54:25Z, 2013-10-18T13:53:18Z](62272852) Equal
2020-12-12T20:23:26.248Z - [2013-10-18T15:53:27Z, 2013-10-18T15:58:38Z](2) MissingInB
2020-12-12T20:23:26.248Z - [2013-10-18T16:06:13Z, 2015-09-28T14:06:52Z](58012719) Equal
- Done Restoring Database

åreal	1479m2.641s
user	358m0.091s
sys	47m30.605s
```

## Restore phase 2 - plan and prep

Final rollup - we will accumulate in postgres.
- Pre phase-2
  - seed with ./data/jsonl-ted-rollup.20150928.100 (result of phase-1 above)
  - restore these over postgres
    -	ted.watt.2016-02-14-1555.sql.bz2 - last watt backup including history from 2008-07-30 00:04:40
    - ted.watt-just2016.2016-02-14-1624.sql.gz - first backup with table truncated from 2016-01-01
    - ted.watt.20201120.2332Z.sql.bz2 - most recent backup with table truncated from 2016-01-01

```go
verify("mysql(watt) <-> postgres", mysql.NewReader(db, "watt"), postgres.NewReader(conn, "watt"))
doTest("mysql(watt) -> postgres", mysql.NewReader(db, "watt"), postgres.NewWriter(conn, "watt"))
verify("mysql(watt) <-> postgres", mysql.NewReader(db, "watt"), postgres.NewReader(conn, "watt"))
```

- Phase-2 Verification
  - restore each `ted.watt*.sql.bz2` and verify against postgres.

```go
verify("mysql(watt) <-> postgres", mysql.NewReader(db, "watt"), postgres.NewReader(conn, "watt"))
```

## Pre phase-2

  - seed with ./data/jsonl-ted-rollup.20150928.100 (result of phase-1 above)
  - restore these over postgres
    -	ted.watt.2016-02-14-1555.sql.bz2 - last watt backup including history from 2008-07-30 00:04:40
    - ted.watt-just2016.2016-02-14-1624.sql.gz - first backup with table truncated from 2016-01-01
    - ted.watt.20201120.2332Z.sql.bz2 - most recent backup with table truncated from 2016-01-01

### jsonl-ted-rollup.20150928.100

```bash
2020-12-13T05:26:15.254Z - Starting TED1K mysql restore
2020-12-13T05:26:15.271Z - Connected to MySQL
2020-12-13T05:26:15.289Z - Connected to Postgres
2020-12-13T05:26:15.343Z - -=- jsonl <-> postgres
2020-12-13T05:27:16.272Z - jsonl <-> postgres (2009-08-21) inner rate ~ 517.6k/s, took 1m0.929s, rate ~ 517.6k/s count: 31536000
2020-12-13T05:28:16.146Z - jsonl <-> postgres (2010-10-16) inner rate ~ 526.7k/s, took 2m0.803s, rate ~ 522.1k/s count: 63072000
2020-12-13T05:29:15.794Z - jsonl <-> postgres (2011-11-11) inner rate ~ 528.7k/s, took 3m0.45s, rate ~ 524.3k/s count: 94608000
2020-12-13T05:30:15.534Z - jsonl <-> postgres (2012-11-15) inner rate ~ 527.9k/s, took 4m0.191s, rate ~ 525.2k/s count: 126144000
2020-12-13T05:31:15.142Z - jsonl <-> postgres (2013-11-22) inner rate ~ 529.1k/s, took 4m59.799s, rate ~ 526.0k/s count: 157680000
2020-12-13T05:32:14.347Z - jsonl <-> postgres (2014-12-15) inner rate ~ 532.7k/s, took 5m59.003s, rate ~ 527.1k/s count: 189216000
2020-12-13T05:32:58.922Z - jsonl <-> postgres took 6m43.578s, rate ~ 527.1k/s count: 212737970
2020-12-13T05:32:58.922Z - Verified jsonl <-> postgres:
2020-12-13T05:32:58.922Z - [2008-07-30T00:04:40Z, 2015-09-28T14:06:52Z](212737970) Equal
```

### ted.watt.2016-02-14-1555.sql.bz2 - last watt backup including history from 2008-07-30 00:04:40
```bash
-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.2016-02-14-1555.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
5208.896s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2016-02-11 12:22:45	223101124

- run mysql restore
2020-12-13T07:01:32.207Z - Starting TED1K mysql restore
2020-12-13T07:01:32.209Z - Connected to MySQL
2020-12-13T07:01:32.218Z - Connected to Postgres
2020-12-13T07:01:32.231Z - -=- mysql(watt) <-> postgres
2020-12-13T07:04:55.918Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 154.8k/s, took 3m23.687s, rate ~ 154.8k/s count: 31536000
2020-12-13T07:08:19.273Z - mysql(watt) <-> postgres (2010-10-16) inner rate ~ 155.1k/s, took 6m47.041s, rate ~ 155.0k/s count: 63072000
2020-12-13T07:11:43.608Z - mysql(watt) <-> postgres (2011-11-11) inner rate ~ 154.3k/s, took 10m11.377s, rate ~ 154.7k/s count: 94608000
2020-12-13T07:15:05.252Z - mysql(watt) <-> postgres (2012-11-15) inner rate ~ 156.4k/s, took 13m33.021s, rate ~ 155.2k/s count: 126144000
2020-12-13T07:18:25.908Z - mysql(watt) <-> postgres (2013-11-22) inner rate ~ 157.2k/s, took 16m53.677s, rate ~ 155.6k/s count: 157680000
2020-12-13T07:21:53.315Z - mysql(watt) <-> postgres (2014-12-15) inner rate ~ 152.1k/s, took 20m21.083s, rate ~ 155.0k/s count: 189216000
2020-12-13T07:24:22.795Z - mysql(watt) <-> postgres took 22m50.564s, rate ~ 155.2k/s count: 212737970
2020-12-13T07:25:28.728Z - Verified mysql(watt) <-> postgres:
2020-12-13T07:25:28.728Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T07:25:28.728Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T07:25:28.728Z - [2009-11-17T03:00:30Z, 2015-09-28T14:06:52Z](173851864) Equal
2020-12-13T07:25:28.728Z - [2015-09-28T14:06:53Z, 2016-02-11T12:22:45Z](10363179) MissingInB
2020-12-13T07:25:28.728Z - -=- mysql(watt) -> postgres
2020-12-13T07:29:29.816Z - mysql(watt) -> postgres (2009-08-21) inner rate ~ 130.8k/s, took 4m1.088s, rate ~ 130.8k/s count: 31536000
2020-12-13T07:33:29.669Z - mysql(watt) -> postgres (2010-10-16) inner rate ~ 131.5k/s, took 8m0.939s, rate ~ 131.1k/s count: 63072000
2020-12-13T07:37:30.992Z - mysql(watt) -> postgres (2011-11-11) inner rate ~ 130.7k/s, took 12m2.263s, rate ~ 131.0k/s count: 94608000
2020-12-13T07:41:31.194Z - mysql(watt) -> postgres (2012-11-15) inner rate ~ 131.3k/s, took 16m2.466s, rate ~ 131.1k/s count: 126144000
2020-12-13T07:45:30.667Z - mysql(watt) -> postgres (2013-11-22) inner rate ~ 131.7k/s, took 20m1.938s, rate ~ 131.2k/s count: 157680000
2020-12-13T07:49:31.066Z - mysql(watt) -> postgres (2014-12-15) inner rate ~ 131.2k/s, took 24m2.337s, rate ~ 131.2k/s count: 189216000
2020-12-13T07:53:20.327Z - mysql(watt) -> postgres (2015-12-31) inner rate ~ 137.6k/s, took 27m51.598s, rate ~ 132.1k/s count: 220752000
2020-12-13T07:53:36.252Z - mysql(watt) -> postgres took 28m7.523s, rate ~ 132.2k/s count: 223101124
2020-12-13T07:53:36.260Z - -=- mysql(watt) <-> postgres
2020-12-13T07:57:01.656Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 153.5k/s, took 3m25.395s, rate ~ 153.5k/s count: 31536000
2020-12-13T08:00:20.627Z - mysql(watt) <-> postgres (2010-10-16) inner rate ~ 158.5k/s, took 6m44.366s, rate ~ 156.0k/s count: 63072000
2020-12-13T08:03:39.739Z - mysql(watt) <-> postgres (2011-11-11) inner rate ~ 158.4k/s, took 10m3.478s, rate ~ 156.8k/s count: 94608000
2020-12-13T08:06:58.445Z - mysql(watt) <-> postgres (2012-11-15) inner rate ~ 158.7k/s, took 13m22.184s, rate ~ 157.3k/s count: 126144000
2020-12-13T08:10:16.929Z - mysql(watt) <-> postgres (2013-11-22) inner rate ~ 158.9k/s, took 16m40.668s, rate ~ 157.6k/s count: 157680000
2020-12-13T08:13:41.346Z - mysql(watt) <-> postgres (2014-12-15) inner rate ~ 154.3k/s, took 20m5.085s, rate ~ 157.0k/s count: 189216000
2020-12-13T08:17:01.149Z - mysql(watt) <-> postgres (2015-12-31) inner rate ~ 157.8k/s, took 23m24.888s, rate ~ 157.1k/s count: 220752000
2020-12-13T08:17:16.088Z - mysql(watt) <-> postgres took 23m39.827s, rate ~ 157.1k/s count: 223101149
2020-12-13T08:17:16.089Z - Verified mysql(watt) <-> postgres:
2020-12-13T08:17:16.089Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T08:17:16.089Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T08:17:16.089Z - [2009-11-17T03:00:30Z, 2016-02-11T12:22:45Z](184215043) Equal
- Done Restoring Database

real	162m44.152s
```

### ted.watt-just2016.2016-02-14-1624.sql.gz - first backup with table truncated from 2016-01-01

note: first recompress `.gz -> .bz2`

```bash
-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt-just2016.2016-02-14-1624.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
52.089s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-01-01 00:00:00	2016-02-14 21:24:21	2318726

- run mysql restore
2020-12-13T08:22:30.372Z - Starting TED1K mysql restore
2020-12-13T08:22:30.374Z - Connected to MySQL
2020-12-13T08:22:30.386Z - Connected to Postgres
2020-12-13T08:22:30.390Z - -=- mysql(watt) <-> postgres
2020-12-13T08:23:03.062Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 965.3k/s, took 32.671s, rate ~ 965.3k/s count: 31536000
2020-12-13T08:23:24.417Z - mysql(watt) <-> postgres (2010-10-16) inner rate ~ 1.5M/s, took 54.026s, rate ~ 1.2M/s count: 63072000
2020-12-13T08:23:45.241Z - mysql(watt) <-> postgres (2011-11-11) inner rate ~ 1.5M/s, took 1m14.85s, rate ~ 1.3M/s count: 94608000
2020-12-13T08:24:05.880Z - mysql(watt) <-> postgres (2012-11-15) inner rate ~ 1.5M/s, took 1m35.489s, rate ~ 1.3M/s count: 126144000
2020-12-13T08:24:27.115Z - mysql(watt) <-> postgres (2013-11-22) inner rate ~ 1.5M/s, took 1m56.724s, rate ~ 1.4M/s count: 157680000
2020-12-13T08:24:48.008Z - mysql(watt) <-> postgres (2014-12-15) inner rate ~ 1.5M/s, took 2m17.617s, rate ~ 1.4M/s count: 189216000
2020-12-13T08:25:08.242Z - mysql(watt) <-> postgres (2015-12-31) inner rate ~ 1.6M/s, took 2m37.851s, rate ~ 1.4M/s count: 220752000
2020-12-13T08:25:32.496Z - mysql(watt) <-> postgres took 3m2.105s, rate ~ 1.2M/s count: 223101149
2020-12-13T08:25:32.499Z - Verified mysql(watt) <-> postgres:
2020-12-13T08:25:32.499Z - [2008-07-30T00:04:40Z, 2015-12-31T23:59:59Z](220783301) MissingInA
2020-12-13T08:25:32.499Z - [2016-01-01T00:00:00Z, 2016-02-11T12:22:45Z](2317848) Equal
2020-12-13T08:25:32.500Z - [2016-02-14T21:09:15Z, 2016-02-14T21:24:21Z](878) MissingInB
2020-12-13T08:25:32.500Z - -=- mysql(watt) -> postgres
2020-12-13T08:26:09.674Z - mysql(watt) -> postgres took 37.174s, rate ~ 62.4k/s count: 2318726
2020-12-13T08:26:09.773Z - -=- mysql(watt) <-> postgres
2020-12-13T08:26:40.715Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.0M/s, took 30.942s, rate ~ 1.0M/s count: 31536000
2020-12-13T08:27:01.992Z - mysql(watt) <-> postgres (2010-10-16) inner rate ~ 1.5M/s, took 52.218s, rate ~ 1.2M/s count: 63072000
2020-12-13T08:27:22.737Z - mysql(watt) <-> postgres (2011-11-11) inner rate ~ 1.5M/s, took 1m12.964s, rate ~ 1.3M/s count: 94608000
2020-12-13T08:27:43.708Z - mysql(watt) <-> postgres (2012-11-15) inner rate ~ 1.5M/s, took 1m33.934s, rate ~ 1.3M/s count: 126144000
2020-12-13T08:28:04.785Z - mysql(watt) <-> postgres (2013-11-22) inner rate ~ 1.5M/s, took 1m55.012s, rate ~ 1.4M/s count: 157680000
2020-12-13T08:28:25.510Z - mysql(watt) <-> postgres (2014-12-15) inner rate ~ 1.5M/s, took 2m15.737s, rate ~ 1.4M/s count: 189216000
2020-12-13T08:28:45.784Z - mysql(watt) <-> postgres (2015-12-31) inner rate ~ 1.6M/s, took 2m36.01s, rate ~ 1.4M/s count: 220752000
2020-12-13T08:29:10.193Z - mysql(watt) <-> postgres took 3m0.42s, rate ~ 1.2M/s count: 223102027
2020-12-13T08:29:10.194Z - Verified mysql(watt) <-> postgres:
2020-12-13T08:29:10.194Z - [2008-07-30T00:04:40Z, 2015-12-31T23:59:59Z](220783301) MissingInA
2020-12-13T08:29:10.194Z - [2016-01-01T00:00:00Z, 2016-02-14T21:24:21Z](2318726) Equal
- Done Restoring Database

real	7m34.955s
```

### ted.watt.20201120.2332Z.sql.bz2 - most recent backup with table truncated from 2016-01-01

```bash
-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20201120.2332Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1986.672s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2020-11-20 23:32:33	143688967

- run mysql restore
2020-12-13T09:03:29.097Z - Starting TED1K mysql restore
2020-12-13T09:03:29.099Z - Connected to MySQL
2020-12-13T09:03:29.108Z - Connected to Postgres
2020-12-13T09:03:29.123Z - -=- mysql(watt) <-> postgres
2020-12-13T09:04:01.578Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 971.7k/s, took 32.455s, rate ~ 971.7k/s count: 31536000
2020-12-13T09:04:26.971Z - mysql(watt) <-> postgres (2010-10-16) inner rate ~ 1.2M/s, took 57.848s, rate ~ 1.1M/s count: 63072000
2020-12-13T09:04:52.264Z - mysql(watt) <-> postgres (2011-11-11) inner rate ~ 1.2M/s, took 1m23.14s, rate ~ 1.1M/s count: 94608000
2020-12-13T09:05:19.569Z - mysql(watt) <-> postgres (2012-11-15) inner rate ~ 1.2M/s, took 1m50.446s, rate ~ 1.1M/s count: 126144000
2020-12-13T09:05:44.822Z - mysql(watt) <-> postgres (2013-11-22) inner rate ~ 1.2M/s, took 2m15.699s, rate ~ 1.2M/s count: 157680000
2020-12-13T09:06:09.642Z - mysql(watt) <-> postgres (2014-12-15) inner rate ~ 1.3M/s, took 2m40.519s, rate ~ 1.2M/s count: 189216000
2020-12-13T09:06:33.778Z - mysql(watt) <-> postgres (2015-12-31) inner rate ~ 1.3M/s, took 3m4.655s, rate ~ 1.2M/s count: 220752000
2020-12-13T09:06:35.319Z - mysql(watt) <-> postgres took 3m6.195s, rate ~ 1.2M/s count: 223102027
2020-12-13T09:10:33.448Z - Verified mysql(watt) <-> postgres:
2020-12-13T09:10:33.449Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInB
2020-12-13T09:10:33.449Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-13T09:10:33.449Z - [2016-03-12T06:35:35Z, 2020-11-20T23:32:33Z](143613941) MissingInB
2020-12-13T09:10:33.449Z - -=- mysql(watt) -> postgres
2020-12-13T09:13:22.411Z - mysql(watt) -> postgres (2017-03-19) inner rate ~ 186.6k/s, took 2m48.962s, rate ~ 186.6k/s count: 31536000
2020-12-13T09:16:10.966Z - mysql(watt) -> postgres (2018-04-16) inner rate ~ 187.1k/s, took 5m37.517s, rate ~ 186.9k/s count: 63072000
2020-12-13T09:18:58.349Z - mysql(watt) -> postgres (2019-04-23) inner rate ~ 188.4k/s, took 8m24.9s, rate ~ 187.4k/s count: 94608000
2020-12-13T09:21:50.477Z - mysql(watt) -> postgres (2020-04-30) inner rate ~ 183.2k/s, took 11m17.028s, rate ~ 186.3k/s count: 126144000
2020-12-13T09:23:34.037Z - mysql(watt) -> postgres took 13m0.588s, rate ~ 184.1k/s count: 143688967
2020-12-13T09:23:34.063Z - -=- mysql(watt) <-> postgres
2020-12-13T09:23:58.303Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 24.24s, rate ~ 1.3M/s count: 31536000
2020-12-13T09:24:25.172Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.2M/s, took 51.108s, rate ~ 1.2M/s count: 63072000
2020-12-13T09:24:49.789Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m15.725s, rate ~ 1.2M/s count: 94608000
2020-12-13T09:25:16.687Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.2M/s, took 1m42.623s, rate ~ 1.2M/s count: 126144000
2020-12-13T09:25:42.576Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.2M/s, took 2m8.512s, rate ~ 1.2M/s count: 157680000
2020-12-13T09:26:06.364Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.3M/s, took 2m32.3s, rate ~ 1.2M/s count: 189216000
2020-12-13T09:26:26.897Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.5M/s, took 2m52.834s, rate ~ 1.3M/s count: 220752000
2020-12-13T09:27:25.248Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 540.5k/s, took 3m51.184s, rate ~ 1.1M/s count: 252288000
2020-12-13T09:28:35.824Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 446.8k/s, took 5m1.761s, rate ~ 940.6k/s count: 283824000
2020-12-13T09:29:46.410Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 446.8k/s, took 6m12.346s, rate ~ 847.0k/s count: 315360000
2020-12-13T09:30:58.094Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 439.9k/s, took 7m24.03s, rate ~ 781.2k/s count: 346896000
2020-12-13T09:31:43.787Z - mysql(watt) <-> postgres took 8m9.724s, rate ~ 749.0k/s count: 366790994
2020-12-13T09:31:43.791Z - Verified mysql(watt) <-> postgres:
2020-12-13T09:31:43.791Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) Equal
2020-12-13T09:31:43.791Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-13T09:31:43.791Z - [2016-03-12T06:35:35Z, 2020-11-20T23:32:33Z](143613941) Equal
- Done Restoring Database
real	61m47.523s
```

## Phase 2 - Verification

```bash
- Verifying docker environment
WARNING: No swap limit support
  Docker seems to be setup properly

- Waiting for database server (go-ted1k_mysql_1) to accept connections (max 30 seconds)
  Connected

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt-just2016.2016-02-14-1624.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
53.968s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-01-01 00:00:00	2016-02-14 21:24:21	2318726

- run mysql restore
2020-12-13T09:36:54.080Z - Starting TED1K mysql restore
2020-12-13T09:36:54.082Z - Connected to MySQL
2020-12-13T09:36:54.092Z - Connected to Postgres
2020-12-13T09:36:54.108Z - -=- mysql(watt) <-> postgres
2020-12-13T09:37:48.588Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 578.9k/s, took 54.479s, rate ~ 578.9k/s count: 31536000
2020-12-13T09:38:16.681Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.1M/s, took 1m22.573s, rate ~ 763.8k/s count: 63072000
2020-12-13T09:38:37.879Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.5M/s, took 1m43.77s, rate ~ 911.7k/s count: 94608000
2020-12-13T09:39:02.092Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.3M/s, took 2m7.983s, rate ~ 985.6k/s count: 126144000
2020-12-13T09:39:28.357Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.2M/s, took 2m34.249s, rate ~ 1.0M/s count: 157680000
2020-12-13T09:39:50.053Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 2m55.945s, rate ~ 1.1M/s count: 189216000
2020-12-13T09:40:12.664Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.4M/s, took 3m18.556s, rate ~ 1.1M/s count: 220752000
2020-12-13T09:41:01.564Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 644.9k/s, took 4m7.455s, rate ~ 1.0M/s count: 252288000
2020-12-13T09:41:23.358Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.4M/s, took 4m29.25s, rate ~ 1.1M/s count: 283824000
2020-12-13T09:41:44.146Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 4m50.038s, rate ~ 1.1M/s count: 315360000
2020-12-13T09:42:08.517Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.3M/s, took 5m14.408s, rate ~ 1.1M/s count: 346896000
2020-12-13T09:42:21.909Z - mysql(watt) <-> postgres took 5m27.8s, rate ~ 1.1M/s count: 366790994
2020-12-13T09:42:21.909Z - Verified mysql(watt) <-> postgres:
2020-12-13T09:42:21.909Z - [2007-08-28T00:02:18Z, 2015-12-31T23:59:59Z](220858327) MissingInA
2020-12-13T09:42:21.909Z - [2016-01-01T00:00:00Z, 2016-02-14T21:24:21Z](2318726) Equal
2020-12-13T09:42:21.909Z - [2016-03-12T06:35:35Z, 2020-11-20T23:32:33Z](143613941) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20090918.0300.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
779.274s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-09-18 07:00:36	33874979

- run mysql restore
2020-12-13T09:55:26.182Z - Starting TED1K mysql restore
2020-12-13T09:55:26.184Z - Connected to MySQL
2020-12-13T09:55:26.192Z - Connected to Postgres
2020-12-13T09:55:26.202Z - -=- mysql(watt) <-> postgres
2020-12-13T09:58:42.099Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 161.0k/s, took 3m15.897s, rate ~ 161.0k/s count: 31536000
2020-12-13T09:59:16.440Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 918.3k/s, took 3m50.237s, rate ~ 273.9k/s count: 63072000
2020-12-13T09:59:36.825Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.5M/s, took 4m10.622s, rate ~ 377.5k/s count: 94608000
2020-12-13T09:59:57.448Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.5M/s, took 4m31.246s, rate ~ 465.1k/s count: 126144000
2020-12-13T10:00:18.101Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.5M/s, took 4m51.898s, rate ~ 540.2k/s count: 157680000
2020-12-13T10:00:38.697Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 5m12.494s, rate ~ 605.5k/s count: 189216000
2020-12-13T10:00:59.320Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.5M/s, took 5m33.117s, rate ~ 662.7k/s count: 220752000
2020-12-13T10:01:20.021Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.5M/s, took 5m53.818s, rate ~ 713.0k/s count: 252288000
2020-12-13T10:01:40.535Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 6m14.333s, rate ~ 758.2k/s count: 283824000
2020-12-13T10:02:01.287Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 6m35.084s, rate ~ 798.2k/s count: 315360000
2020-12-13T10:02:21.852Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.5M/s, took 6m55.649s, rate ~ 834.6k/s count: 346896000
2020-12-13T10:02:34.989Z - mysql(watt) <-> postgres took 7m8.787s, rate ~ 855.4k/s count: 366790994
2020-12-13T10:02:34.990Z - Verified mysql(watt) <-> postgres:
2020-12-13T10:02:34.990Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T10:02:34.990Z - [2008-07-30T00:04:40Z, 2009-09-18T07:00:36Z](33874979) Equal
2020-12-13T10:02:34.990Z - [2009-09-18T07:00:37Z, 2020-11-20T23:32:33Z](332840989) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20091022.0258.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
843.550s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-10-22 06:58:36	36782768

- run mysql restore
2020-12-13T10:16:43.515Z - Starting TED1K mysql restore
2020-12-13T10:16:43.517Z - Connected to MySQL
2020-12-13T10:16:43.529Z - Connected to Postgres
2020-12-13T10:16:43.547Z - -=- mysql(watt) <-> postgres
2020-12-13T10:19:55.615Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 164.2k/s, took 3m12.068s, rate ~ 164.2k/s count: 31536000
2020-12-13T10:20:45.204Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 636.0k/s, took 4m1.657s, rate ~ 261.0k/s count: 63072000
2020-12-13T10:21:05.888Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.5M/s, took 4m22.341s, rate ~ 360.6k/s count: 94608000
2020-12-13T10:21:26.643Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.5M/s, took 4m43.096s, rate ~ 445.6k/s count: 126144000
2020-12-13T10:21:47.191Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.5M/s, took 5m3.644s, rate ~ 519.3k/s count: 157680000
2020-12-13T10:22:07.672Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 5m24.125s, rate ~ 583.8k/s count: 189216000
2020-12-13T10:22:27.787Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.6M/s, took 5m44.24s, rate ~ 641.3k/s count: 220752000
2020-12-13T10:22:48.215Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.5M/s, took 6m4.668s, rate ~ 691.8k/s count: 252288000
2020-12-13T10:23:08.652Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 6m25.105s, rate ~ 737.0k/s count: 283824000
2020-12-13T10:23:29.071Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 6m45.524s, rate ~ 777.7k/s count: 315360000
2020-12-13T10:23:49.519Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.5M/s, took 7m5.972s, rate ~ 814.4k/s count: 346896000
2020-12-13T10:24:02.744Z - mysql(watt) <-> postgres took 7m19.197s, rate ~ 835.1k/s count: 366790994
2020-12-13T10:24:02.745Z - Verified mysql(watt) <-> postgres:
2020-12-13T10:24:02.745Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T10:24:02.745Z - [2008-07-30T00:04:40Z, 2009-10-22T06:58:36Z](36782768) Equal
2020-12-13T10:24:02.745Z - [2009-10-22T06:58:37Z, 2020-11-20T23:32:33Z](329933200) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20091102.0134.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
856.466s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2009-11-02 06:34:50	37626188

- run mysql restore
2020-12-13T10:38:24.126Z - Starting TED1K mysql restore
2020-12-13T10:38:24.128Z - Connected to MySQL
2020-12-13T10:38:24.143Z - Connected to Postgres
2020-12-13T10:38:24.161Z - -=- mysql(watt) <-> postgres
2020-12-13T10:41:33.980Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 166.1k/s, took 3m9.819s, rate ~ 166.1k/s count: 31536000
2020-12-13T10:42:27.643Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 587.7k/s, took 4m3.482s, rate ~ 259.0k/s count: 63072000
2020-12-13T10:42:48.632Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.5M/s, took 4m24.471s, rate ~ 357.7k/s count: 94608000
2020-12-13T10:43:09.191Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.5M/s, took 4m45.03s, rate ~ 442.6k/s count: 126144000
2020-12-13T10:43:29.839Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.5M/s, took 5m5.678s, rate ~ 515.8k/s count: 157680000
2020-12-13T10:43:50.295Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 5m26.134s, rate ~ 580.2k/s count: 189216000
2020-12-13T10:44:10.885Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.5M/s, took 5m46.724s, rate ~ 636.7k/s count: 220752000
2020-12-13T10:44:31.423Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.5M/s, took 6m7.262s, rate ~ 686.9k/s count: 252288000
2020-12-13T10:44:51.836Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 6m27.675s, rate ~ 732.1k/s count: 283824000
2020-12-13T10:45:12.488Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 6m48.328s, rate ~ 772.3k/s count: 315360000
2020-12-13T10:45:32.838Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.5M/s, took 7m8.677s, rate ~ 809.2k/s count: 346896000
2020-12-13T10:45:45.629Z - mysql(watt) <-> postgres took 7m21.468s, rate ~ 830.8k/s count: 366790994
2020-12-13T10:45:45.630Z - Verified mysql(watt) <-> postgres:
2020-12-13T10:45:45.630Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T10:45:45.630Z - [2008-07-30T00:04:40Z, 2009-11-02T06:34:50Z](37626188) Equal
2020-12-13T10:45:45.630Z - [2009-11-02T06:34:51Z, 2020-11-20T23:32:33Z](329089780) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20110406.0316.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1799.550s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2011-04-06 07:15:35	77646482

- run mysql restore
2020-12-13T11:15:48.897Z - Starting TED1K mysql restore
2020-12-13T11:15:48.899Z - Connected to MySQL
2020-12-13T11:15:48.912Z - Connected to Postgres
2020-12-13T11:15:48.933Z - -=- mysql(watt) <-> postgres
2020-12-13T11:19:10.544Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 156.4k/s, took 3m21.611s, rate ~ 156.4k/s count: 31536000
2020-12-13T11:22:30.682Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 157.6k/s, took 6m41.749s, rate ~ 157.0k/s count: 63072000
2020-12-13T11:24:14.446Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 303.9k/s, took 8m25.513s, rate ~ 187.2k/s count: 94608000
2020-12-13T11:24:35.082Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.5M/s, took 8m46.149s, rate ~ 239.7k/s count: 126144000
2020-12-13T11:24:55.629Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.5M/s, took 9m6.696s, rate ~ 288.4k/s count: 157680000
2020-12-13T11:25:16.614Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 9m27.681s, rate ~ 333.3k/s count: 189216000
2020-12-13T11:25:37.710Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.5M/s, took 9m48.777s, rate ~ 374.9k/s count: 220752000
2020-12-13T11:25:58.361Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.5M/s, took 10m9.428s, rate ~ 414.0k/s count: 252288000
2020-12-13T11:26:19.158Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 10m30.225s, rate ~ 450.4k/s count: 283824000
2020-12-13T11:26:39.826Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 10m50.893s, rate ~ 484.5k/s count: 315360000
2020-12-13T11:27:00.106Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.6M/s, took 11m11.173s, rate ~ 516.9k/s count: 346896000
2020-12-13T11:27:13.017Z - mysql(watt) <-> postgres took 11m24.085s, rate ~ 536.2k/s count: 366790994
2020-12-13T11:27:13.018Z - Verified mysql(watt) <-> postgres:
2020-12-13T11:27:13.018Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T11:27:13.018Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T11:27:13.018Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T11:27:13.018Z - [2009-11-17T03:00:30Z, 2011-04-06T07:15:35Z](38760401) Equal
2020-12-13T11:27:13.018Z - [2011-04-06T07:24:52Z, 2020-11-20T23:32:33Z](289069461) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20110607.0115.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1917.007s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2011-06-07 05:15:59	82931473

- run mysql restore
2020-12-13T11:59:15.113Z - Starting TED1K mysql restore
2020-12-13T11:59:15.115Z - Connected to MySQL
2020-12-13T11:59:15.124Z - Connected to Postgres
2020-12-13T11:59:15.142Z - -=- mysql(watt) <-> postgres
2020-12-13T12:02:30.585Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 161.4k/s, took 3m15.443s, rate ~ 161.4k/s count: 31536000
2020-12-13T12:05:45.891Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 161.5k/s, took 6m30.749s, rate ~ 161.4k/s count: 63072000
2020-12-13T12:08:00.199Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 234.8k/s, took 8m45.057s, rate ~ 180.2k/s count: 94608000
2020-12-13T12:08:28.029Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.1M/s, took 9m12.887s, rate ~ 228.2k/s count: 126144000
2020-12-13T12:08:57.157Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.1M/s, took 9m42.015s, rate ~ 270.9k/s count: 157680000
2020-12-13T12:09:23.544Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.2M/s, took 10m8.402s, rate ~ 311.0k/s count: 189216000
2020-12-13T12:09:49.903Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.2M/s, took 10m34.761s, rate ~ 347.8k/s count: 220752000
2020-12-13T12:10:19.164Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.1M/s, took 11m4.022s, rate ~ 379.9k/s count: 252288000
2020-12-13T12:10:45.897Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.2M/s, took 11m30.755s, rate ~ 410.9k/s count: 283824000
2020-12-13T12:11:12.418Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.2M/s, took 11m57.276s, rate ~ 439.7k/s count: 315360000
2020-12-13T12:11:39.453Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.2M/s, took 12m24.311s, rate ~ 466.1k/s count: 346896000
2020-12-13T12:11:54.486Z - mysql(watt) <-> postgres took 12m39.344s, rate ~ 483.0k/s count: 366790994
2020-12-13T12:11:54.488Z - Verified mysql(watt) <-> postgres:
2020-12-13T12:11:54.488Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T12:11:54.488Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T12:11:54.488Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T12:11:54.488Z - [2009-11-17T03:00:30Z, 2011-06-07T05:15:59Z](44045392) Equal
2020-12-13T12:11:54.488Z - [2011-06-07T05:16:00Z, 2020-11-20T23:32:33Z](283784470) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20111017.2045.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
2122.491s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2011-10-18 00:33:13	92547284

- run mysql restore
2020-12-13T12:47:21.973Z - Starting TED1K mysql restore
2020-12-13T12:47:21.975Z - Connected to MySQL
2020-12-13T12:47:21.982Z - Connected to Postgres
2020-12-13T12:47:22.003Z - -=- mysql(watt) <-> postgres
2020-12-13T12:50:41.565Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 158.0k/s, took 3m19.562s, rate ~ 158.0k/s count: 31536000
2020-12-13T12:54:00.897Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 158.2k/s, took 6m38.894s, rate ~ 158.1k/s count: 63072000
2020-12-13T12:57:09.256Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 167.4k/s, took 9m47.254s, rate ~ 161.1k/s count: 94608000
2020-12-13T12:57:29.927Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.5M/s, took 10m7.924s, rate ~ 207.5k/s count: 126144000
2020-12-13T12:57:50.404Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.5M/s, took 10m28.401s, rate ~ 250.9k/s count: 157680000
2020-12-13T12:58:10.939Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 10m48.937s, rate ~ 291.6k/s count: 189216000
2020-12-13T12:58:31.685Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.5M/s, took 11m9.682s, rate ~ 329.6k/s count: 220752000
2020-12-13T12:58:52.331Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.5M/s, took 11m30.328s, rate ~ 365.5k/s count: 252288000
2020-12-13T12:59:12.884Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 11m50.881s, rate ~ 399.3k/s count: 283824000
2020-12-13T12:59:33.567Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 12m11.565s, rate ~ 431.1k/s count: 315360000
2020-12-13T12:59:54.191Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.5M/s, took 12m32.188s, rate ~ 461.2k/s count: 346896000
2020-12-13T13:00:07.160Z - mysql(watt) <-> postgres took 12m45.157s, rate ~ 479.4k/s count: 366790994
2020-12-13T13:00:07.160Z - Verified mysql(watt) <-> postgres:
2020-12-13T13:00:07.160Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T13:00:07.160Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T13:00:07.160Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T13:00:07.160Z - [2009-11-17T03:00:30Z, 2011-10-18T00:33:13Z](53661203) Equal
2020-12-13T13:00:07.160Z - [2011-10-18T00:51:14Z, 2020-11-20T23:32:33Z](274168659) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20120608.0119.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
2584.931s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2012-06-08 05:19:20	112499188

- run mysql restore
2020-12-13T13:43:17.091Z - Starting TED1K mysql restore
2020-12-13T13:43:17.093Z - Connected to MySQL
2020-12-13T13:43:17.100Z - Connected to Postgres
2020-12-13T13:43:17.132Z - -=- mysql(watt) <-> postgres
2020-12-13T13:46:37.183Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 157.6k/s, took 3m20.051s, rate ~ 157.6k/s count: 31536000
2020-12-13T13:49:55.471Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 159.0k/s, took 6m38.338s, rate ~ 158.3k/s count: 63072000
2020-12-13T13:53:15.040Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 158.0k/s, took 9m57.908s, rate ~ 158.2k/s count: 94608000
2020-12-13T13:55:17.189Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 258.2k/s, took 12m0.057s, rate ~ 175.2k/s count: 126144000
2020-12-13T13:55:38.297Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.5M/s, took 12m21.165s, rate ~ 212.7k/s count: 157680000
2020-12-13T13:55:59.142Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 12m42.01s, rate ~ 248.3k/s count: 189216000
2020-12-13T13:56:19.848Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.5M/s, took 13m2.716s, rate ~ 282.0k/s count: 220752000
2020-12-13T13:56:40.420Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.5M/s, took 13m23.288s, rate ~ 314.1k/s count: 252288000
2020-12-13T13:57:01.341Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 13m44.208s, rate ~ 344.4k/s count: 283824000
2020-12-13T13:57:22.146Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 14m5.014s, rate ~ 373.2k/s count: 315360000
2020-12-13T13:57:42.779Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.5M/s, took 14m25.647s, rate ~ 400.7k/s count: 346896000
2020-12-13T13:57:55.849Z - mysql(watt) <-> postgres took 14m38.717s, rate ~ 417.4k/s count: 366790994
2020-12-13T13:57:55.850Z - Verified mysql(watt) <-> postgres:
2020-12-13T13:57:55.850Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T13:57:55.851Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T13:57:55.851Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T13:57:55.851Z - [2009-11-17T03:00:30Z, 2012-06-08T05:19:20Z](73613107) Equal
2020-12-13T13:57:55.851Z - [2012-06-08T05:19:22Z, 2020-11-20T23:32:33Z](254216755) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20130221.2119.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
3089.150s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2013-02-22 02:19:09	134540091

- run mysql restore
2020-12-13T14:49:29.985Z - Starting TED1K mysql restore
2020-12-13T14:49:29.987Z - Connected to MySQL
2020-12-13T14:49:29.994Z - Connected to Postgres
2020-12-13T14:49:30.015Z - -=- mysql(watt) <-> postgres
2020-12-13T14:52:49.827Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 157.8k/s, took 3m19.812s, rate ~ 157.8k/s count: 31536000
2020-12-13T14:56:09.983Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 157.6k/s, took 6m39.968s, rate ~ 157.7k/s count: 63072000
2020-12-13T14:59:29.261Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 158.3k/s, took 9m59.246s, rate ~ 157.9k/s count: 94608000
2020-12-13T15:02:49.172Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 157.8k/s, took 13m19.157s, rate ~ 157.8k/s count: 126144000
2020-12-13T15:03:57.885Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 459.0k/s, took 14m27.87s, rate ~ 181.7k/s count: 157680000
2020-12-13T15:04:18.696Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 14m48.68s, rate ~ 212.9k/s count: 189216000
2020-12-13T15:04:39.285Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.5M/s, took 15m9.27s, rate ~ 242.8k/s count: 220752000
2020-12-13T15:05:00.134Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.5M/s, took 15m30.119s, rate ~ 271.2k/s count: 252288000
2020-12-13T15:05:20.550Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 15m50.535s, rate ~ 298.6k/s count: 283824000
2020-12-13T15:05:41.029Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 16m11.014s, rate ~ 324.8k/s count: 315360000
2020-12-13T15:06:01.730Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.5M/s, took 16m31.715s, rate ~ 349.8k/s count: 346896000
2020-12-13T15:06:14.637Z - mysql(watt) <-> postgres took 16m44.622s, rate ~ 365.1k/s count: 366790994
2020-12-13T15:06:14.637Z - Verified mysql(watt) <-> postgres:
2020-12-13T15:06:14.637Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T15:06:14.637Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T15:06:14.637Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T15:06:14.638Z - [2009-11-17T03:00:30Z, 2013-02-22T02:19:09Z](95654010) Equal
2020-12-13T15:06:14.638Z - [2013-02-22T02:19:10Z, 2020-11-20T23:32:33Z](232175852) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20140219.2038.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
3812.441s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2014-02-19 13:58:18	165309499

- run mysql restore
2020-12-13T16:09:51.356Z - Starting TED1K mysql restore
2020-12-13T16:09:51.358Z - Connected to MySQL
2020-12-13T16:09:51.365Z - Connected to Postgres
2020-12-13T16:09:51.386Z - -=- mysql(watt) <-> postgres
2020-12-13T16:13:14.378Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 155.4k/s, took 3m22.992s, rate ~ 155.4k/s count: 31536000
2020-12-13T16:16:41.367Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 152.4k/s, took 6m49.981s, rate ~ 153.8k/s count: 63072000
2020-12-13T16:20:25.189Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 140.9k/s, took 10m33.801s, rate ~ 149.3k/s count: 94608000
2020-12-13T16:24:09.893Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 140.3k/s, took 14m18.506s, rate ~ 146.9k/s count: 126144000
2020-12-13T16:27:52.568Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 141.6k/s, took 18m1.181s, rate ~ 145.8k/s count: 157680000
2020-12-13T16:29:18.079Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 368.8k/s, took 19m26.693s, rate ~ 162.2k/s count: 189216000
2020-12-13T16:29:51.708Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 937.8k/s, took 20m0.322s, rate ~ 183.9k/s count: 220752000
2020-12-13T16:30:24.092Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 973.8k/s, took 20m32.705s, rate ~ 204.7k/s count: 252288000
2020-12-13T16:30:58.748Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 909.9k/s, took 21m7.362s, rate ~ 223.9k/s count: 283824000
2020-12-13T16:31:30.038Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.0M/s, took 21m38.651s, rate ~ 242.8k/s count: 315360000
2020-12-13T16:31:59.244Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.1M/s, took 22m7.857s, rate ~ 261.2k/s count: 346896000
2020-12-13T16:32:18.842Z - mysql(watt) <-> postgres took 22m27.456s, rate ~ 272.2k/s count: 366790994
2020-12-13T16:32:18.844Z - Verified mysql(watt) <-> postgres:
2020-12-13T16:32:18.844Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T16:32:18.844Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T16:32:18.845Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T16:32:18.845Z - [2009-11-17T03:00:30Z, 2014-02-19T13:58:18Z](126423418) Equal
2020-12-13T16:32:18.845Z - [2014-02-20T01:44:19Z, 2020-11-20T23:32:33Z](201406444) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20140806.0016.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
4679.499s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2014-08-06 04:16:19	179484943

- run mysql restore
2020-12-13T17:50:26.541Z - Starting TED1K mysql restore
2020-12-13T17:50:26.543Z - Connected to MySQL
2020-12-13T17:50:26.559Z - Connected to Postgres
2020-12-13T17:50:26.610Z - -=- mysql(watt) <-> postgres
2020-12-13T17:55:07.999Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 112.1k/s, took 4m41.389s, rate ~ 112.1k/s count: 31536000
2020-12-13T17:59:32.926Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 119.0k/s, took 9m6.315s, rate ~ 115.4k/s count: 63072000
2020-12-13T18:03:12.100Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 143.9k/s, took 12m45.489s, rate ~ 123.6k/s count: 94608000
2020-12-13T18:06:59.775Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 138.5k/s, took 16m33.165s, rate ~ 127.0k/s count: 126144000
2020-12-13T18:10:44.192Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 140.5k/s, took 20m17.582s, rate ~ 129.5k/s count: 157680000
2020-12-13T18:13:20.587Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 201.6k/s, took 22m53.976s, rate ~ 137.7k/s count: 189216000
2020-12-13T18:13:44.499Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 23m17.889s, rate ~ 157.9k/s count: 220752000
2020-12-13T18:14:09.346Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.3M/s, took 23m42.736s, rate ~ 177.3k/s count: 252288000
2020-12-13T18:14:35.361Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.2M/s, took 24m8.751s, rate ~ 195.9k/s count: 283824000
2020-12-13T18:15:03.366Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.1M/s, took 24m36.756s, rate ~ 213.5k/s count: 315360000
2020-12-13T18:15:27.982Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.3M/s, took 25m1.372s, rate ~ 231.1k/s count: 346896000
2020-12-13T18:15:48.315Z - mysql(watt) <-> postgres took 25m21.705s, rate ~ 241.0k/s count: 366790994
2020-12-13T18:15:48.316Z - Verified mysql(watt) <-> postgres:
2020-12-13T18:15:48.316Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T18:15:48.316Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T18:15:48.316Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T18:15:48.316Z - [2009-11-17T03:00:30Z, 2014-08-06T04:16:19Z](140598862) Equal
2020-12-13T18:15:48.316Z - [2014-08-06T04:16:20Z, 2020-11-20T23:32:33Z](187231000) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20141005.2218.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
4462.593s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2014-10-06 02:18:37	184124340

- run mysql restore
2020-12-13T19:30:16.474Z - Starting TED1K mysql restore
2020-12-13T19:30:16.476Z - Connected to MySQL
2020-12-13T19:30:16.492Z - Connected to Postgres
2020-12-13T19:30:16.547Z - -=- mysql(watt) <-> postgres
2020-12-13T19:33:44.069Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 152.0k/s, took 3m27.522s, rate ~ 152.0k/s count: 31536000
2020-12-13T19:37:08.114Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 154.6k/s, took 6m51.567s, rate ~ 153.2k/s count: 63072000
2020-12-13T19:40:33.124Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 153.8k/s, took 10m16.578s, rate ~ 153.4k/s count: 94608000
2020-12-13T19:43:56.008Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 155.4k/s, took 13m39.461s, rate ~ 153.9k/s count: 126144000
2020-12-13T19:47:20.154Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 154.5k/s, took 17m3.607s, rate ~ 154.0k/s count: 157680000
2020-12-13T19:50:14.608Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 180.8k/s, took 19m58.061s, rate ~ 157.9k/s count: 189216000
2020-12-13T19:50:36.155Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.5M/s, took 20m19.609s, rate ~ 181.0k/s count: 220752000
2020-12-13T19:50:57.795Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.5M/s, took 20m41.249s, rate ~ 203.3k/s count: 252288000
2020-12-13T19:51:19.473Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 21m2.926s, rate ~ 224.7k/s count: 283824000
2020-12-13T19:51:41.225Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.4M/s, took 21m24.679s, rate ~ 245.5k/s count: 315360000
2020-12-13T19:52:02.764Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.5M/s, took 21m46.218s, rate ~ 265.6k/s count: 346896000
2020-12-13T19:52:16.738Z - mysql(watt) <-> postgres took 22m0.191s, rate ~ 277.8k/s count: 366790994
2020-12-13T19:52:16.739Z - Verified mysql(watt) <-> postgres:
2020-12-13T19:52:16.739Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T19:52:16.739Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T19:52:16.739Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T19:52:16.739Z - [2009-11-17T03:00:30Z, 2014-10-06T02:18:37Z](145238259) Equal
2020-12-13T19:52:16.739Z - [2014-10-06T02:18:38Z, 2020-11-20T23:32:33Z](182591603) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20150928.1003.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
5166.468s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2015-09-28 14:03:18	212737731

- run mysql restore
2020-12-13T21:18:28.507Z - Starting TED1K mysql restore
2020-12-13T21:18:28.510Z - Connected to MySQL
2020-12-13T21:18:28.525Z - Connected to Postgres
2020-12-13T21:18:28.539Z - -=- mysql(watt) <-> postgres
2020-12-13T21:21:53.660Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 153.7k/s, took 3m25.121s, rate ~ 153.7k/s count: 31536000
2020-12-13T21:25:15.369Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 156.3k/s, took 6m46.83s, rate ~ 155.0k/s count: 63072000
2020-12-13T21:28:37.590Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 155.9k/s, took 10m9.051s, rate ~ 155.3k/s count: 94608000
2020-12-13T21:32:00.109Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 155.7k/s, took 13m31.57s, rate ~ 155.4k/s count: 126144000
2020-12-13T21:35:21.018Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 157.0k/s, took 16m52.478s, rate ~ 155.7k/s count: 157680000
2020-12-13T21:38:42.028Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 156.9k/s, took 20m13.488s, rate ~ 155.9k/s count: 189216000
2020-12-13T21:41:18.123Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 202.0k/s, took 22m49.584s, rate ~ 161.2k/s count: 220752000
2020-12-13T21:41:39.451Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 1.5M/s, took 23m10.912s, rate ~ 181.4k/s count: 252288000
2020-12-13T21:42:00.556Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 23m32.017s, rate ~ 201.0k/s count: 283824000
2020-12-13T21:42:22.020Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 23m53.481s, rate ~ 220.0k/s count: 315360000
2020-12-13T21:42:43.567Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.5M/s, took 24m15.028s, rate ~ 238.4k/s count: 346896000
2020-12-13T21:42:57.265Z - mysql(watt) <-> postgres took 24m28.726s, rate ~ 249.7k/s count: 366790994
2020-12-13T21:42:57.266Z - Verified mysql(watt) <-> postgres:
2020-12-13T21:42:57.266Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T21:42:57.266Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T21:42:57.266Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T21:42:57.266Z - [2009-11-17T03:00:30Z, 2015-09-28T14:03:18Z](173851650) Equal
2020-12-13T21:42:57.266Z - [2015-09-28T14:03:19Z, 2020-11-20T23:32:33Z](153978212) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.2016-02-14-1555.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
5407.622s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2008-07-30 00:04:40	2016-02-11 12:22:45	223101124

- run mysql restore
2020-12-13T23:13:10.082Z - Starting TED1K mysql restore
2020-12-13T23:13:10.084Z - Connected to MySQL
2020-12-13T23:13:10.098Z - Connected to Postgres
2020-12-13T23:13:10.126Z - -=- mysql(watt) <-> postgres
2020-12-13T23:16:36.586Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 152.7k/s, took 3m26.459s, rate ~ 152.7k/s count: 31536000
2020-12-13T23:20:01.925Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 153.6k/s, took 6m51.799s, rate ~ 153.2k/s count: 63072000
2020-12-13T23:23:27.723Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 153.2k/s, took 10m17.596s, rate ~ 153.2k/s count: 94608000
2020-12-13T23:26:50.066Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 155.9k/s, took 13m39.939s, rate ~ 153.8k/s count: 126144000
2020-12-13T23:30:13.402Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 155.1k/s, took 17m3.275s, rate ~ 154.1k/s count: 157680000
2020-12-13T23:33:37.223Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 154.7k/s, took 20m27.096s, rate ~ 154.2k/s count: 189216000
2020-12-13T23:37:00.710Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 155.0k/s, took 23m50.583s, rate ~ 154.3k/s count: 220752000
2020-12-13T23:37:38.341Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 838.0k/s, took 24m28.215s, rate ~ 171.8k/s count: 252288000
2020-12-13T23:38:03.611Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.2M/s, took 24m53.484s, rate ~ 190.0k/s count: 283824000
2020-12-13T23:38:28.807Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.3M/s, took 25m18.68s, rate ~ 207.7k/s count: 315360000
2020-12-13T23:38:55.816Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.2M/s, took 25m45.69s, rate ~ 224.4k/s count: 346896000
2020-12-13T23:39:12.030Z - mysql(watt) <-> postgres took 26m1.903s, rate ~ 234.8k/s count: 366790994
2020-12-13T23:39:12.031Z - Verified mysql(watt) <-> postgres:
2020-12-13T23:39:12.031Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInA
2020-12-13T23:39:12.031Z - [2008-07-30T00:04:40Z, 2009-11-17T02:46:02Z](38886081) Equal
2020-12-13T23:39:12.031Z - [2009-11-17T02:46:03Z, 2009-11-17T02:46:27Z](25) MissingInA
2020-12-13T23:39:12.031Z - [2009-11-17T03:00:30Z, 2016-02-11T12:22:45Z](184215043) Equal
2020-12-13T23:39:12.031Z - [2016-02-14T21:09:15Z, 2020-11-20T23:32:33Z](143614819) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20160430.0232Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
56.878s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2016-04-30 02:32:57	4047909

- run mysql restore
2020-12-13T23:40:14.905Z - Starting TED1K mysql restore
2020-12-13T23:40:14.907Z - Connected to MySQL
2020-12-13T23:40:14.919Z - Connected to Postgres
2020-12-13T23:40:14.938Z - -=- mysql(watt) <-> postgres
2020-12-13T23:40:38.690Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 23.751s, rate ~ 1.3M/s count: 31536000
2020-12-13T23:41:02.977Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 48.038s, rate ~ 1.3M/s count: 63072000
2020-12-13T23:41:26.631Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m11.692s, rate ~ 1.3M/s count: 94608000
2020-12-13T23:41:52.249Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.2M/s, took 1m37.31s, rate ~ 1.3M/s count: 126144000
2020-12-13T23:42:17.682Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.2M/s, took 2m2.743s, rate ~ 1.3M/s count: 157680000
2020-12-13T23:42:47.691Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.1M/s, took 2m32.753s, rate ~ 1.2M/s count: 189216000
2020-12-13T23:43:12.456Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 2m57.517s, rate ~ 1.2M/s count: 220752000
2020-12-13T23:43:45.648Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 950.1k/s, took 3m30.71s, rate ~ 1.2M/s count: 252288000
2020-12-13T23:44:14.106Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.1M/s, took 3m59.168s, rate ~ 1.2M/s count: 283824000
2020-12-13T23:44:43.017Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.1M/s, took 4m28.078s, rate ~ 1.2M/s count: 315360000
2020-12-13T23:45:13.988Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.0M/s, took 4m59.049s, rate ~ 1.2M/s count: 346896000
2020-12-13T23:45:32.964Z - mysql(watt) <-> postgres took 5m18.025s, rate ~ 1.2M/s count: 366790994
2020-12-13T23:45:32.964Z - Verified mysql(watt) <-> postgres:
2020-12-13T23:45:32.964Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-13T23:45:32.964Z - [2016-03-12T06:35:35Z, 2016-04-30T02:32:57Z](4047909) Equal
2020-12-13T23:45:32.964Z - [2016-04-30T02:32:58Z, 2020-11-20T23:32:33Z](139566032) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20160616.0229Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
111.123s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2016-06-16 02:29:18	8067692

- run mysql restore
2020-12-13T23:47:30.201Z - Starting TED1K mysql restore
2020-12-13T23:47:30.203Z - Connected to MySQL
2020-12-13T23:47:30.212Z - Connected to Postgres
2020-12-13T23:47:30.230Z - -=- mysql(watt) <-> postgres
2020-12-13T23:47:57.942Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.1M/s, took 27.711s, rate ~ 1.1M/s count: 31536000
2020-12-13T23:48:21.044Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.4M/s, took 50.814s, rate ~ 1.2M/s count: 63072000
2020-12-13T23:48:43.837Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.4M/s, took 1m13.607s, rate ~ 1.3M/s count: 94608000
2020-12-13T23:49:10.950Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.2M/s, took 1m40.72s, rate ~ 1.3M/s count: 126144000
2020-12-13T23:49:39.917Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.1M/s, took 2m9.687s, rate ~ 1.2M/s count: 157680000
2020-12-13T23:50:06.606Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.2M/s, took 2m36.376s, rate ~ 1.2M/s count: 189216000
2020-12-13T23:50:33.698Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.2M/s, took 3m3.467s, rate ~ 1.2M/s count: 220752000
2020-12-13T23:51:09.696Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 876.0k/s, took 3m39.466s, rate ~ 1.1M/s count: 252288000
2020-12-13T23:51:39.545Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.1M/s, took 4m9.315s, rate ~ 1.1M/s count: 283824000
2020-12-13T23:52:04.442Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.3M/s, took 4m34.211s, rate ~ 1.2M/s count: 315360000
2020-12-13T23:52:30.343Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.2M/s, took 5m0.113s, rate ~ 1.2M/s count: 346896000
2020-12-13T23:52:46.730Z - mysql(watt) <-> postgres took 5m16.499s, rate ~ 1.2M/s count: 366790994
2020-12-13T23:52:46.730Z - Verified mysql(watt) <-> postgres:
2020-12-13T23:52:46.730Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-13T23:52:46.730Z - [2016-03-12T06:35:35Z, 2016-06-16T02:29:18Z](8067692) Equal
2020-12-13T23:52:46.730Z - [2016-06-16T02:29:19Z, 2020-11-20T23:32:33Z](135546249) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20160719.1848Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
150.445s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2016-07-19 18:48:44	10855735

- run mysql restore
2020-12-13T23:55:23.508Z - Starting TED1K mysql restore
2020-12-13T23:55:23.510Z - Connected to MySQL
2020-12-13T23:55:23.520Z - Connected to Postgres
2020-12-13T23:55:23.540Z - -=- mysql(watt) <-> postgres
2020-12-13T23:55:49.925Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.2M/s, took 26.385s, rate ~ 1.2M/s count: 31536000
2020-12-13T23:56:16.491Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.2M/s, took 52.95s, rate ~ 1.2M/s count: 63072000
2020-12-13T23:56:41.482Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m17.941s, rate ~ 1.2M/s count: 94608000
2020-12-13T23:57:08.918Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.1M/s, took 1m45.377s, rate ~ 1.2M/s count: 126144000
2020-12-13T23:57:35.848Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.2M/s, took 2m12.307s, rate ~ 1.2M/s count: 157680000
2020-12-13T23:58:00.890Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.3M/s, took 2m37.349s, rate ~ 1.2M/s count: 189216000
2020-12-13T23:58:23.811Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.4M/s, took 3m0.27s, rate ~ 1.2M/s count: 220752000
2020-12-13T23:59:01.386Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 839.3k/s, took 3m37.845s, rate ~ 1.2M/s count: 252288000
2020-12-13T23:59:22.745Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.5M/s, took 3m59.205s, rate ~ 1.2M/s count: 283824000
2020-12-13T23:59:44.426Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.5M/s, took 4m20.886s, rate ~ 1.2M/s count: 315360000
2020-12-14T00:00:06.341Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.4M/s, took 4m42.8s, rate ~ 1.2M/s count: 346896000
2020-12-14T00:00:21.190Z - mysql(watt) <-> postgres took 4m57.649s, rate ~ 1.2M/s count: 366790994
2020-12-14T00:00:21.190Z - Verified mysql(watt) <-> postgres:
2020-12-14T00:00:21.190Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T00:00:21.190Z - [2016-03-12T06:35:35Z, 2016-07-19T18:48:44Z](10855735) Equal
2020-12-14T00:00:21.190Z - [2016-07-19T18:48:45Z, 2020-11-20T23:32:33Z](132758206) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20160918.0059Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
224.704s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2016-09-18 00:59:58	16004523

- run mysql restore
2020-12-14T00:04:12.629Z - Starting TED1K mysql restore
2020-12-14T00:04:12.631Z - Connected to MySQL
2020-12-14T00:04:12.640Z - Connected to Postgres
2020-12-14T00:04:12.660Z - -=- mysql(watt) <-> postgres
2020-12-14T00:04:39.370Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.2M/s, took 26.71s, rate ~ 1.2M/s count: 31536000
2020-12-14T00:05:03.733Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 51.072s, rate ~ 1.2M/s count: 63072000
2020-12-14T00:05:26.117Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.4M/s, took 1m13.457s, rate ~ 1.3M/s count: 94608000
2020-12-14T00:05:52.572Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.2M/s, took 1m39.912s, rate ~ 1.3M/s count: 126144000
2020-12-14T00:06:14.649Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.4M/s, took 2m1.988s, rate ~ 1.3M/s count: 157680000
2020-12-14T00:06:39.102Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.3M/s, took 2m26.442s, rate ~ 1.3M/s count: 189216000
2020-12-14T00:07:02.414Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.4M/s, took 2m49.754s, rate ~ 1.3M/s count: 220752000
2020-12-14T00:07:46.019Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 723.2k/s, took 3m33.359s, rate ~ 1.2M/s count: 252288000
2020-12-14T00:08:11.628Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.2M/s, took 3m58.967s, rate ~ 1.2M/s count: 283824000
2020-12-14T00:08:35.690Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.3M/s, took 4m23.03s, rate ~ 1.2M/s count: 315360000
2020-12-14T00:09:00.132Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.3M/s, took 4m47.471s, rate ~ 1.2M/s count: 346896000
2020-12-14T00:09:14.721Z - mysql(watt) <-> postgres took 5m2.061s, rate ~ 1.2M/s count: 366790994
2020-12-14T00:09:14.722Z - Verified mysql(watt) <-> postgres:
2020-12-14T00:09:14.722Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T00:09:14.722Z - [2016-03-12T06:35:35Z, 2016-09-18T00:59:58Z](16004523) Equal
2020-12-14T00:09:14.722Z - [2016-09-18T00:59:59Z, 2020-11-20T23:32:33Z](127609418) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20161202.0733Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
318.949s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2016-12-02 07:33:54	22421521

- run mysql restore
2020-12-14T00:14:40.960Z - Starting TED1K mysql restore
2020-12-14T00:14:40.961Z - Connected to MySQL
2020-12-14T00:14:40.970Z - Connected to Postgres
2020-12-14T00:14:40.990Z - -=- mysql(watt) <-> postgres
2020-12-14T00:15:05.170Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 24.18s, rate ~ 1.3M/s count: 31536000
2020-12-14T00:15:28.115Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.4M/s, took 47.125s, rate ~ 1.3M/s count: 63072000
2020-12-14T00:15:49.785Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.5M/s, took 1m8.795s, rate ~ 1.4M/s count: 94608000
2020-12-14T00:16:13.215Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.3M/s, took 1m32.225s, rate ~ 1.4M/s count: 126144000
2020-12-14T00:16:37.687Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.3M/s, took 1m56.697s, rate ~ 1.4M/s count: 157680000
2020-12-14T00:17:03.170Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.2M/s, took 2m22.18s, rate ~ 1.3M/s count: 189216000
2020-12-14T00:17:27.000Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 2m46.01s, rate ~ 1.3M/s count: 220752000
2020-12-14T00:18:20.016Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 594.8k/s, took 3m39.026s, rate ~ 1.2M/s count: 252288000
2020-12-14T00:18:43.476Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.3M/s, took 4m2.486s, rate ~ 1.2M/s count: 283824000
2020-12-14T00:19:06.865Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.3M/s, took 4m25.875s, rate ~ 1.2M/s count: 315360000
2020-12-14T00:19:29.841Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.4M/s, took 4m48.851s, rate ~ 1.2M/s count: 346896000
2020-12-14T00:19:45.667Z - mysql(watt) <-> postgres took 5m4.677s, rate ~ 1.2M/s count: 366790994
2020-12-14T00:19:45.667Z - Verified mysql(watt) <-> postgres:
2020-12-14T00:19:45.667Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T00:19:45.668Z - [2016-03-12T06:35:35Z, 2016-12-02T07:33:54Z](22421521) Equal
2020-12-14T00:19:45.668Z - [2016-12-02T07:33:55Z, 2020-11-20T23:32:33Z](121192420) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20170106.0629Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
361.681s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2017-01-06 06:29:54	25369345

- run mysql restore
2020-12-14T00:25:54.903Z - Starting TED1K mysql restore
2020-12-14T00:25:54.906Z - Connected to MySQL
2020-12-14T00:25:54.915Z - Connected to Postgres
2020-12-14T00:25:54.936Z - -=- mysql(watt) <-> postgres
2020-12-14T00:26:23.995Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.1M/s, took 29.058s, rate ~ 1.1M/s count: 31536000
2020-12-14T00:26:47.423Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 52.487s, rate ~ 1.2M/s count: 63072000
2020-12-14T00:27:09.103Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.5M/s, took 1m14.166s, rate ~ 1.3M/s count: 94608000
2020-12-14T00:27:32.717Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.3M/s, took 1m37.781s, rate ~ 1.3M/s count: 126144000
2020-12-14T00:27:55.792Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.4M/s, took 2m0.855s, rate ~ 1.3M/s count: 157680000
2020-12-14T00:28:18.468Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.4M/s, took 2m23.532s, rate ~ 1.3M/s count: 189216000
2020-12-14T00:28:43.067Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 2m48.13s, rate ~ 1.3M/s count: 220752000
2020-12-14T00:29:40.781Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 546.4k/s, took 3m45.844s, rate ~ 1.1M/s count: 252288000
2020-12-14T00:30:06.792Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.2M/s, took 4m11.856s, rate ~ 1.1M/s count: 283824000
2020-12-14T00:30:36.034Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.1M/s, took 4m41.097s, rate ~ 1.1M/s count: 315360000
2020-12-14T00:31:04.124Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.1M/s, took 5m9.187s, rate ~ 1.1M/s count: 346896000
2020-12-14T00:31:17.852Z - mysql(watt) <-> postgres took 5m22.916s, rate ~ 1.1M/s count: 366790994
2020-12-14T00:31:17.853Z - Verified mysql(watt) <-> postgres:
2020-12-14T00:31:17.853Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T00:31:17.853Z - [2016-03-12T06:35:35Z, 2017-01-06T06:29:54Z](25369345) Equal
2020-12-14T00:31:17.853Z - [2017-01-06T06:29:55Z, 2020-11-20T23:32:33Z](118244596) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20170326.1528Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
456.608s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2017-03-26 15:28:49	32081147

- run mysql restore
2020-12-14T00:39:02.599Z - Starting TED1K mysql restore
2020-12-14T00:39:02.601Z - Connected to MySQL
2020-12-14T00:39:02.610Z - Connected to Postgres
2020-12-14T00:39:02.629Z - -=- mysql(watt) <-> postgres
2020-12-14T00:39:27.318Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 24.689s, rate ~ 1.3M/s count: 31536000
2020-12-14T00:39:52.300Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 49.671s, rate ~ 1.3M/s count: 63072000
2020-12-14T00:40:16.869Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m14.24s, rate ~ 1.3M/s count: 94608000
2020-12-14T00:40:39.828Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.4M/s, took 1m37.198s, rate ~ 1.3M/s count: 126144000
2020-12-14T00:41:02.704Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.4M/s, took 2m0.075s, rate ~ 1.3M/s count: 157680000
2020-12-14T00:41:28.372Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.2M/s, took 2m25.743s, rate ~ 1.3M/s count: 189216000
2020-12-14T00:41:54.708Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.2M/s, took 2m52.079s, rate ~ 1.3M/s count: 220752000
2020-12-14T00:42:57.655Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 501.0k/s, took 3m55.026s, rate ~ 1.1M/s count: 252288000
2020-12-14T00:43:26.114Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.1M/s, took 4m23.484s, rate ~ 1.1M/s count: 283824000
2020-12-14T00:43:50.563Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.3M/s, took 4m47.934s, rate ~ 1.1M/s count: 315360000
2020-12-14T00:44:14.264Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.3M/s, took 5m11.635s, rate ~ 1.1M/s count: 346896000
2020-12-14T00:44:28.737Z - mysql(watt) <-> postgres took 5m26.108s, rate ~ 1.1M/s count: 366790994
2020-12-14T00:44:28.738Z - Verified mysql(watt) <-> postgres:
2020-12-14T00:44:28.738Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T00:44:28.738Z - [2016-03-12T06:35:35Z, 2017-03-26T15:28:49Z](32081147) Equal
2020-12-14T00:44:28.738Z - [2017-03-26T15:28:50Z, 2020-11-20T23:32:33Z](111532794) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20170607.0541Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
541.331s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2017-06-07 05:41:49	38228226

- run mysql restore
2020-12-14T00:53:38.736Z - Starting TED1K mysql restore
2020-12-14T00:53:38.738Z - Connected to MySQL
2020-12-14T00:53:38.747Z - Connected to Postgres
2020-12-14T00:53:38.768Z - -=- mysql(watt) <-> postgres
2020-12-14T00:54:06.892Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.1M/s, took 28.124s, rate ~ 1.1M/s count: 31536000
2020-12-14T00:54:32.436Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.2M/s, took 53.667s, rate ~ 1.2M/s count: 63072000
2020-12-14T00:54:56.329Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m17.56s, rate ~ 1.2M/s count: 94608000
2020-12-14T00:55:27.504Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.0M/s, took 1m48.736s, rate ~ 1.2M/s count: 126144000
2020-12-14T00:55:53.685Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.2M/s, took 2m14.917s, rate ~ 1.2M/s count: 157680000
2020-12-14T00:56:16.521Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.4M/s, took 2m37.753s, rate ~ 1.2M/s count: 189216000
2020-12-14T00:56:40.077Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 3m1.309s, rate ~ 1.2M/s count: 220752000
2020-12-14T00:57:43.378Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 498.2k/s, took 4m4.609s, rate ~ 1.0M/s count: 252288000
2020-12-14T00:58:19.757Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 866.9k/s, took 4m40.989s, rate ~ 1.0M/s count: 283824000
2020-12-14T00:58:44.954Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.3M/s, took 5m6.186s, rate ~ 1.0M/s count: 315360000
2020-12-14T00:59:09.005Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.3M/s, took 5m30.236s, rate ~ 1.1M/s count: 346896000
2020-12-14T00:59:23.874Z - mysql(watt) <-> postgres took 5m45.105s, rate ~ 1.1M/s count: 366790994
2020-12-14T00:59:23.874Z - Verified mysql(watt) <-> postgres:
2020-12-14T00:59:23.874Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T00:59:23.874Z - [2016-03-12T06:35:35Z, 2017-06-07T05:41:49Z](38228226) Equal
2020-12-14T00:59:23.874Z - [2017-06-07T05:41:50Z, 2020-11-20T23:32:33Z](105385715) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20170727.1724Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
597.669s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2017-07-27 17:24:30	41955521

- run mysql restore
2020-12-14T01:09:30.464Z - Starting TED1K mysql restore
2020-12-14T01:09:30.466Z - Connected to MySQL
2020-12-14T01:09:30.475Z - Connected to Postgres
2020-12-14T01:09:30.499Z - -=- mysql(watt) <-> postgres
2020-12-14T01:09:54.267Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 23.768s, rate ~ 1.3M/s count: 31536000
2020-12-14T01:10:16.420Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.4M/s, took 45.921s, rate ~ 1.4M/s count: 63072000
2020-12-14T01:10:39.027Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.4M/s, took 1m8.528s, rate ~ 1.4M/s count: 94608000
2020-12-14T01:11:02.110Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.4M/s, took 1m31.611s, rate ~ 1.4M/s count: 126144000
2020-12-14T01:11:26.338Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.3M/s, took 1m55.839s, rate ~ 1.4M/s count: 157680000
2020-12-14T01:11:49.101Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.4M/s, took 2m18.602s, rate ~ 1.4M/s count: 189216000
2020-12-14T01:12:11.513Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.4M/s, took 2m41.014s, rate ~ 1.4M/s count: 220752000
2020-12-14T01:13:14.082Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 504.0k/s, took 3m43.583s, rate ~ 1.1M/s count: 252288000
2020-12-14T01:13:54.242Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 785.3k/s, took 4m23.743s, rate ~ 1.1M/s count: 283824000
2020-12-14T01:14:17.830Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.3M/s, took 4m47.331s, rate ~ 1.1M/s count: 315360000
2020-12-14T01:14:44.087Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.2M/s, took 5m13.589s, rate ~ 1.1M/s count: 346896000
2020-12-14T01:14:59.868Z - mysql(watt) <-> postgres took 5m29.369s, rate ~ 1.1M/s count: 366790994
2020-12-14T01:14:59.868Z - Verified mysql(watt) <-> postgres:
2020-12-14T01:14:59.868Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T01:14:59.869Z - [2016-03-12T06:35:35Z, 2017-07-27T17:24:30Z](41955521) Equal
2020-12-14T01:14:59.869Z - [2017-07-27T17:24:31Z, 2020-11-20T23:32:33Z](101658420) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20180217.2219Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
817.176s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2018-02-17 22:19:46	58328447

- run mysql restore
2020-12-14T01:28:48.326Z - Starting TED1K mysql restore
2020-12-14T01:28:48.327Z - Connected to MySQL
2020-12-14T01:28:48.337Z - Connected to Postgres
2020-12-14T01:28:48.358Z - -=- mysql(watt) <-> postgres
2020-12-14T01:29:12.475Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 24.117s, rate ~ 1.3M/s count: 31536000
2020-12-14T01:29:35.844Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 47.486s, rate ~ 1.3M/s count: 63072000
2020-12-14T01:29:58.981Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.4M/s, took 1m10.623s, rate ~ 1.3M/s count: 94608000
2020-12-14T01:30:21.561Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.4M/s, took 1m33.203s, rate ~ 1.4M/s count: 126144000
2020-12-14T01:30:45.554Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.3M/s, took 1m57.197s, rate ~ 1.3M/s count: 157680000
2020-12-14T01:31:11.885Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.2M/s, took 2m23.527s, rate ~ 1.3M/s count: 189216000
2020-12-14T01:31:37.043Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 2m48.685s, rate ~ 1.3M/s count: 220752000
2020-12-14T01:32:39.982Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 501.1k/s, took 3m51.624s, rate ~ 1.1M/s count: 252288000
2020-12-14T01:33:41.355Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 513.8k/s, took 4m52.997s, rate ~ 968.7k/s count: 283824000
2020-12-14T01:34:07.440Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.2M/s, took 5m19.082s, rate ~ 988.3k/s count: 315360000
2020-12-14T01:34:30.806Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.3M/s, took 5m42.449s, rate ~ 1.0M/s count: 346896000
2020-12-14T01:34:46.188Z - mysql(watt) <-> postgres took 5m57.83s, rate ~ 1.0M/s count: 366790994
2020-12-14T01:34:46.189Z - Verified mysql(watt) <-> postgres:
2020-12-14T01:34:46.189Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T01:34:46.189Z - [2016-03-12T06:35:35Z, 2018-02-17T22:19:46Z](58328447) Equal
2020-12-14T01:34:46.189Z - [2018-02-17T22:19:47Z, 2020-11-20T23:32:33Z](85285494) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20180326.0312Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
852.732s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2018-03-26 03:12:25	61205052

- run mysql restore
2020-12-14T01:49:09.542Z - Starting TED1K mysql restore
2020-12-14T01:49:09.544Z - Connected to MySQL
2020-12-14T01:49:09.553Z - Connected to Postgres
2020-12-14T01:49:09.574Z - -=- mysql(watt) <-> postgres
2020-12-14T01:49:32.615Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.4M/s, took 23.042s, rate ~ 1.4M/s count: 31536000
2020-12-14T01:49:58.005Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.2M/s, took 48.431s, rate ~ 1.3M/s count: 63072000
2020-12-14T01:50:22.353Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m12.779s, rate ~ 1.3M/s count: 94608000
2020-12-14T01:50:48.757Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.2M/s, took 1m39.183s, rate ~ 1.3M/s count: 126144000
2020-12-14T01:51:17.770Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.1M/s, took 2m8.196s, rate ~ 1.2M/s count: 157680000
2020-12-14T01:51:43.275Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.2M/s, took 2m33.701s, rate ~ 1.2M/s count: 189216000
2020-12-14T01:52:09.046Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.2M/s, took 2m59.472s, rate ~ 1.2M/s count: 220752000
2020-12-14T01:53:14.346Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 482.9k/s, took 4m4.772s, rate ~ 1.0M/s count: 252288000
2020-12-14T01:54:20.922Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 473.7k/s, took 5m11.348s, rate ~ 911.6k/s count: 283824000
2020-12-14T01:54:50.030Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.1M/s, took 5m40.456s, rate ~ 926.3k/s count: 315360000
2020-12-14T01:55:16.458Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.2M/s, took 6m6.884s, rate ~ 945.5k/s count: 346896000
2020-12-14T01:55:31.967Z - mysql(watt) <-> postgres took 6m22.393s, rate ~ 959.2k/s count: 366790994
2020-12-14T01:55:31.967Z - Verified mysql(watt) <-> postgres:
2020-12-14T01:55:31.967Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T01:55:31.967Z - [2016-03-12T06:35:35Z, 2018-03-26T03:12:25Z](61205052) Equal
2020-12-14T01:55:31.967Z - [2018-03-26T03:12:26Z, 2020-11-20T23:32:33Z](82408889) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20180612.0035Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
937.995s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2018-06-12 00:35:24	67871341

- run mysql restore
2020-12-14T02:11:21.302Z - Starting TED1K mysql restore
2020-12-14T02:11:21.304Z - Connected to MySQL
2020-12-14T02:11:21.313Z - Connected to Postgres
2020-12-14T02:11:21.332Z - -=- mysql(watt) <-> postgres
2020-12-14T02:11:46.859Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.2M/s, took 25.527s, rate ~ 1.2M/s count: 31536000
2020-12-14T02:12:12.129Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.2M/s, took 50.797s, rate ~ 1.2M/s count: 63072000
2020-12-14T02:12:37.277Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m15.945s, rate ~ 1.2M/s count: 94608000
2020-12-14T02:13:04.788Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.1M/s, took 1m43.456s, rate ~ 1.2M/s count: 126144000
2020-12-14T02:13:31.207Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.2M/s, took 2m9.874s, rate ~ 1.2M/s count: 157680000
2020-12-14T02:13:55.324Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.3M/s, took 2m33.992s, rate ~ 1.2M/s count: 189216000
2020-12-14T02:14:19.864Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 2m58.532s, rate ~ 1.2M/s count: 220752000
2020-12-14T02:15:21.476Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 511.8k/s, took 4m0.144s, rate ~ 1.1M/s count: 252288000
2020-12-14T02:16:26.416Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 485.6k/s, took 5m5.083s, rate ~ 930.3k/s count: 283824000
2020-12-14T02:17:04.214Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 834.3k/s, took 5m42.881s, rate ~ 919.7k/s count: 315360000
2020-12-14T02:17:28.270Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.3M/s, took 6m6.938s, rate ~ 945.4k/s count: 346896000
2020-12-14T02:17:46.275Z - mysql(watt) <-> postgres took 6m24.942s, rate ~ 952.8k/s count: 366790994
2020-12-14T02:17:46.275Z - Verified mysql(watt) <-> postgres:
2020-12-14T02:17:46.275Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T02:17:46.275Z - [2016-03-12T06:35:35Z, 2018-06-12T00:35:24Z](67871341) Equal
2020-12-14T02:17:46.275Z - [2018-06-12T00:35:25Z, 2020-11-20T23:32:33Z](75742600) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20180720.2138Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
986.129s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2018-07-20 21:38:28	71229301

- run mysql restore
2020-12-14T02:34:23.744Z - Starting TED1K mysql restore
2020-12-14T02:34:23.747Z - Connected to MySQL
2020-12-14T02:34:23.759Z - Connected to Postgres
2020-12-14T02:34:23.779Z - -=- mysql(watt) <-> postgres
2020-12-14T02:34:50.128Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.2M/s, took 26.349s, rate ~ 1.2M/s count: 31536000
2020-12-14T02:35:14.167Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 50.388s, rate ~ 1.3M/s count: 63072000
2020-12-14T02:35:40.206Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.2M/s, took 1m16.427s, rate ~ 1.2M/s count: 94608000
2020-12-14T02:36:02.547Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.4M/s, took 1m38.768s, rate ~ 1.3M/s count: 126144000
2020-12-14T02:36:26.886Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.3M/s, took 2m3.107s, rate ~ 1.3M/s count: 157680000
2020-12-14T02:36:52.592Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.2M/s, took 2m28.813s, rate ~ 1.3M/s count: 189216000
2020-12-14T02:37:17.572Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 2m53.793s, rate ~ 1.3M/s count: 220752000
2020-12-14T02:38:20.690Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 499.6k/s, took 3m56.911s, rate ~ 1.1M/s count: 252288000
2020-12-14T02:39:25.808Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 484.3k/s, took 5m2.03s, rate ~ 939.7k/s count: 283824000
2020-12-14T02:40:03.220Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 842.9k/s, took 5m39.441s, rate ~ 929.1k/s count: 315360000
2020-12-14T02:40:25.851Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.4M/s, took 6m2.072s, rate ~ 958.1k/s count: 346896000
2020-12-14T02:40:39.977Z - mysql(watt) <-> postgres took 6m16.198s, rate ~ 975.0k/s count: 366790994
2020-12-14T02:40:39.977Z - Verified mysql(watt) <-> postgres:
2020-12-14T02:40:39.977Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T02:40:39.977Z - [2016-03-12T06:35:35Z, 2018-07-20T21:38:28Z](71229301) Equal
2020-12-14T02:40:39.977Z - [2018-07-20T21:38:29Z, 2020-11-20T23:32:33Z](72384640) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20180831.2033Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1044.196s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2016-03-12 06:35:35	2018-08-31 20:33:16	74852305

- run mysql restore
2020-12-14T02:58:17.977Z - Starting TED1K mysql restore
2020-12-14T02:58:17.979Z - Connected to MySQL
2020-12-14T02:58:17.990Z - Connected to Postgres
2020-12-14T02:58:18.012Z - -=- mysql(watt) <-> postgres
2020-12-14T02:58:43.817Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.2M/s, took 25.805s, rate ~ 1.2M/s count: 31536000
2020-12-14T02:59:07.904Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 49.892s, rate ~ 1.3M/s count: 63072000
2020-12-14T02:59:31.771Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m13.759s, rate ~ 1.3M/s count: 94608000
2020-12-14T02:59:53.218Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.5M/s, took 1m35.206s, rate ~ 1.3M/s count: 126144000
2020-12-14T03:00:15.238Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.4M/s, took 1m57.226s, rate ~ 1.3M/s count: 157680000
2020-12-14T03:00:36.937Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 2m18.925s, rate ~ 1.4M/s count: 189216000
2020-12-14T03:01:02.999Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.2M/s, took 2m44.987s, rate ~ 1.3M/s count: 220752000
2020-12-14T03:02:03.144Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 524.3k/s, took 3m45.132s, rate ~ 1.1M/s count: 252288000
2020-12-14T03:03:08.965Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 479.1k/s, took 4m50.952s, rate ~ 975.5k/s count: 283824000
2020-12-14T03:03:49.810Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 772.1k/s, took 5m31.797s, rate ~ 950.5k/s count: 315360000
2020-12-14T03:04:15.500Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.2M/s, took 5m57.488s, rate ~ 970.4k/s count: 346896000
2020-12-14T03:04:32.219Z - mysql(watt) <-> postgres took 6m14.206s, rate ~ 980.2k/s count: 366790994
2020-12-14T03:04:32.219Z - Verified mysql(watt) <-> postgres:
2020-12-14T03:04:32.219Z - [2007-08-28T00:02:18Z, 2016-02-14T21:24:21Z](223177053) MissingInA
2020-12-14T03:04:32.219Z - [2016-03-12T06:35:35Z, 2018-08-31T20:33:16Z](74852305) Equal
2020-12-14T03:04:32.219Z - [2018-08-31T20:33:17Z, 2020-11-20T23:32:33Z](68761636) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20181024.1913Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1137.277s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2018-10-24 19:13:53	79223872

- run mysql restore
2020-12-14T03:23:44.484Z - Starting TED1K mysql restore
2020-12-14T03:23:44.486Z - Connected to MySQL
2020-12-14T03:23:44.495Z - Connected to Postgres
2020-12-14T03:23:44.516Z - -=- mysql(watt) <-> postgres
2020-12-14T03:24:09.301Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 24.785s, rate ~ 1.3M/s count: 31536000
2020-12-14T03:24:32.305Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.4M/s, took 47.789s, rate ~ 1.3M/s count: 63072000
2020-12-14T03:24:58.724Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.2M/s, took 1m14.208s, rate ~ 1.3M/s count: 94608000
2020-12-14T03:25:23.380Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.3M/s, took 1m38.863s, rate ~ 1.3M/s count: 126144000
2020-12-14T03:25:48.131Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.3M/s, took 2m3.614s, rate ~ 1.3M/s count: 157680000
2020-12-14T03:26:15.404Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.2M/s, took 2m30.888s, rate ~ 1.3M/s count: 189216000
2020-12-14T03:26:40.614Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 2m56.097s, rate ~ 1.3M/s count: 220752000
2020-12-14T03:27:41.787Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 515.5k/s, took 3m57.271s, rate ~ 1.1M/s count: 252288000
2020-12-14T03:28:46.979Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 483.7k/s, took 5m2.463s, rate ~ 938.4k/s count: 283824000
2020-12-14T03:29:33.680Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 675.3k/s, took 5m49.164s, rate ~ 903.2k/s count: 315360000
2020-12-14T03:29:57.024Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.4M/s, took 6m12.508s, rate ~ 931.2k/s count: 346896000
2020-12-14T03:30:14.764Z - mysql(watt) <-> postgres took 6m30.247s, rate ~ 939.9k/s count: 366790994
2020-12-14T03:30:14.764Z - Verified mysql(watt) <-> postgres:
2020-12-14T03:30:14.764Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) Equal
2020-12-14T03:30:14.764Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-14T03:30:14.764Z - [2016-03-12T06:35:35Z, 2018-10-24T19:13:53Z](79148846) Equal
2020-12-14T03:30:14.764Z - [2018-10-24T19:13:54Z, 2020-11-20T23:32:33Z](64465095) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20190414.0128Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
ERROR 1064 (42000) at line 920: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near '^[[1' at line 1
read unix @->/var/run/docker.sock: read: connection reset by peer
434.041s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2017-03-16 23:52:32	31336872

- run mysql restore
2020-12-14T03:37:37.378Z - Starting TED1K mysql restore
2020-12-14T03:37:37.380Z - Connected to MySQL
2020-12-14T03:37:37.389Z - Connected to Postgres
2020-12-14T03:37:37.436Z - -=- mysql(watt) <-> postgres
2020-12-14T03:38:03.471Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.2M/s, took 26.035s, rate ~ 1.2M/s count: 31536000
2020-12-14T03:38:28.129Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 50.693s, rate ~ 1.2M/s count: 63072000
2020-12-14T03:38:52.442Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m15.005s, rate ~ 1.3M/s count: 94608000
2020-12-14T03:39:19.076Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.2M/s, took 1m41.64s, rate ~ 1.2M/s count: 126144000
2020-12-14T03:39:43.460Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.3M/s, took 2m6.024s, rate ~ 1.3M/s count: 157680000
2020-12-14T03:40:05.174Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.5M/s, took 2m27.737s, rate ~ 1.3M/s count: 189216000
2020-12-14T03:40:28.266Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.4M/s, took 2m50.83s, rate ~ 1.3M/s count: 220752000
2020-12-14T03:41:30.681Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 505.3k/s, took 3m53.245s, rate ~ 1.1M/s count: 252288000
2020-12-14T03:41:56.238Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 1.2M/s, took 4m18.802s, rate ~ 1.1M/s count: 283824000
2020-12-14T03:42:19.643Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 1.3M/s, took 4m42.206s, rate ~ 1.1M/s count: 315360000
2020-12-14T03:42:41.997Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 1.4M/s, took 5m4.561s, rate ~ 1.1M/s count: 346896000
2020-12-14T03:42:56.669Z - mysql(watt) <-> postgres took 5m19.233s, rate ~ 1.1M/s count: 366790994
2020-12-14T03:42:56.669Z - Verified mysql(watt) <-> postgres:
2020-12-14T03:42:56.669Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) Equal
2020-12-14T03:42:56.669Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-14T03:42:56.669Z - [2016-03-12T06:35:35Z, 2017-03-16T23:52:32Z](31261846) Equal
2020-12-14T03:42:56.669Z - [2017-03-16T23:52:33Z, 2020-11-20T23:32:33Z](112352095) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20190617.0443Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1422.637s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2019-06-17 04:43:57	99352625

- run mysql restore
2020-12-14T04:06:56.221Z - Starting TED1K mysql restore
2020-12-14T04:06:56.224Z - Connected to MySQL
2020-12-14T04:06:56.234Z - Connected to Postgres
2020-12-14T04:06:56.288Z - -=- mysql(watt) <-> postgres
2020-12-14T04:07:19.986Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 23.698s, rate ~ 1.3M/s count: 31536000
2020-12-14T04:07:41.525Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.5M/s, took 45.237s, rate ~ 1.4M/s count: 63072000
2020-12-14T04:08:03.077Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.5M/s, took 1m6.789s, rate ~ 1.4M/s count: 94608000
2020-12-14T04:08:24.742Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.5M/s, took 1m28.454s, rate ~ 1.4M/s count: 126144000
2020-12-14T04:08:46.577Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.4M/s, took 1m50.289s, rate ~ 1.4M/s count: 157680000
2020-12-14T04:09:11.034Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.3M/s, took 2m14.745s, rate ~ 1.4M/s count: 189216000
2020-12-14T04:09:37.981Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.2M/s, took 2m41.692s, rate ~ 1.4M/s count: 220752000
2020-12-14T04:10:40.233Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 506.6k/s, took 3m43.945s, rate ~ 1.1M/s count: 252288000
2020-12-14T04:11:46.410Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 476.5k/s, took 4m50.122s, rate ~ 978.3k/s count: 283824000
2020-12-14T04:12:52.192Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 479.4k/s, took 5m55.904s, rate ~ 886.1k/s count: 315360000
2020-12-14T04:13:25.123Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 957.6k/s, took 6m28.835s, rate ~ 892.1k/s count: 346896000
2020-12-14T04:13:43.991Z - mysql(watt) <-> postgres took 6m47.703s, rate ~ 899.7k/s count: 366790994
2020-12-14T04:13:43.992Z - Verified mysql(watt) <-> postgres:
2020-12-14T04:13:43.992Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) Equal
2020-12-14T04:13:43.992Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-14T04:13:43.992Z - [2016-03-12T06:35:35Z, 2019-06-17T04:43:57Z](99277599) Equal
2020-12-14T04:13:43.992Z - [2019-06-17T04:43:58Z, 2020-11-20T23:32:33Z](44336342) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20190818.0554Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1467.114s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2019-08-18 05:54:32	104330337

- run mysql restore
2020-12-14T04:38:32.909Z - Starting TED1K mysql restore
2020-12-14T04:38:32.911Z - Connected to MySQL
2020-12-14T04:38:32.920Z - Connected to Postgres
2020-12-14T04:38:32.970Z - -=- mysql(watt) <-> postgres
2020-12-14T04:38:57.793Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 24.823s, rate ~ 1.3M/s count: 31536000
2020-12-14T04:39:21.684Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 48.714s, rate ~ 1.3M/s count: 63072000
2020-12-14T04:39:45.185Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m12.215s, rate ~ 1.3M/s count: 94608000
2020-12-14T04:40:07.256Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.4M/s, took 1m34.286s, rate ~ 1.3M/s count: 126144000
2020-12-14T04:40:32.345Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.3M/s, took 1m59.375s, rate ~ 1.3M/s count: 157680000
2020-12-14T04:40:58.601Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.2M/s, took 2m25.631s, rate ~ 1.3M/s count: 189216000
2020-12-14T04:41:23.816Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.3M/s, took 2m50.846s, rate ~ 1.3M/s count: 220752000
2020-12-14T04:42:31.376Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 466.8k/s, took 3m58.406s, rate ~ 1.1M/s count: 252288000
2020-12-14T04:43:39.719Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 461.4k/s, took 5m6.749s, rate ~ 925.3k/s count: 283824000
2020-12-14T04:44:46.878Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 469.6k/s, took 6m13.908s, rate ~ 843.4k/s count: 315360000
2020-12-14T04:45:31.391Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 708.5k/s, took 6m58.421s, rate ~ 829.1k/s count: 346896000
2020-12-14T04:45:48.737Z - mysql(watt) <-> postgres took 7m15.768s, rate ~ 841.7k/s count: 366790994
2020-12-14T04:45:48.739Z - Verified mysql(watt) <-> postgres:
2020-12-14T04:45:48.739Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) Equal
2020-12-14T04:45:48.739Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-14T04:45:48.739Z - [2016-03-12T06:35:35Z, 2019-08-18T05:54:32Z](104255311) Equal
2020-12-14T04:45:48.739Z - [2019-08-20T21:47:49Z, 2020-11-20T23:32:33Z](39358630) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20191129.0710Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1592.204s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2019-11-28 16:13:05	112928156

- run mysql restore
2020-12-14T05:12:43.767Z - Starting TED1K mysql restore
2020-12-14T05:12:43.768Z - Connected to MySQL
2020-12-14T05:12:43.777Z - Connected to Postgres
2020-12-14T05:12:43.827Z - -=- mysql(watt) <-> postgres
2020-12-14T05:13:11.518Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.1M/s, took 27.69s, rate ~ 1.1M/s count: 31536000
2020-12-14T05:13:36.085Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.3M/s, took 52.257s, rate ~ 1.2M/s count: 63072000
2020-12-14T05:14:00.634Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.3M/s, took 1m16.806s, rate ~ 1.2M/s count: 94608000
2020-12-14T05:14:27.816Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.2M/s, took 1m43.989s, rate ~ 1.2M/s count: 126144000
2020-12-14T05:14:54.656Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.2M/s, took 2m10.829s, rate ~ 1.2M/s count: 157680000
2020-12-14T05:15:22.151Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.1M/s, took 2m38.324s, rate ~ 1.2M/s count: 189216000
2020-12-14T05:15:50.141Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.1M/s, took 3m6.313s, rate ~ 1.2M/s count: 220752000
2020-12-14T05:17:04.089Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 426.5k/s, took 4m20.261s, rate ~ 969.4k/s count: 252288000
2020-12-14T05:18:22.671Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 401.3k/s, took 5m38.843s, rate ~ 837.6k/s count: 283824000
2020-12-14T05:19:39.955Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 408.1k/s, took 6m56.127s, rate ~ 757.8k/s count: 315360000
2020-12-14T05:20:34.049Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 583.0k/s, took 7m50.221s, rate ~ 737.7k/s count: 346896000
2020-12-14T05:20:49.433Z - mysql(watt) <-> postgres took 8m5.605s, rate ~ 755.3k/s count: 366790994
2020-12-14T05:20:49.434Z - Verified mysql(watt) <-> postgres:
2020-12-14T05:20:49.434Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) Equal
2020-12-14T05:20:49.434Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-14T05:20:49.434Z - [2016-03-12T06:35:35Z, 2019-11-28T16:13:05Z](112853130) Equal
2020-12-14T05:20:49.434Z - [2019-11-29T07:17:08Z, 2020-11-20T23:32:33Z](30760811) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20200413.1503Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1756.273s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2020-04-13 15:03:10	124702892

- run mysql restore
2020-12-14T05:50:27.858Z - Starting TED1K mysql restore
2020-12-14T05:50:27.860Z - Connected to MySQL
2020-12-14T05:50:27.868Z - Connected to Postgres
2020-12-14T05:50:27.915Z - -=- mysql(watt) <-> postgres
2020-12-14T05:50:51.948Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.3M/s, took 24.032s, rate ~ 1.3M/s count: 31536000
2020-12-14T05:51:13.476Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.5M/s, took 45.56s, rate ~ 1.4M/s count: 63072000
2020-12-14T05:51:35.856Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.4M/s, took 1m7.94s, rate ~ 1.4M/s count: 94608000
2020-12-14T05:51:58.616Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.4M/s, took 1m30.701s, rate ~ 1.4M/s count: 126144000
2020-12-14T05:52:21.393Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.4M/s, took 1m53.478s, rate ~ 1.4M/s count: 157680000
2020-12-14T05:52:43.906Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.4M/s, took 2m15.991s, rate ~ 1.4M/s count: 189216000
2020-12-14T05:53:06.716Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.4M/s, took 2m38.8s, rate ~ 1.4M/s count: 220752000
2020-12-14T05:54:08.428Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 511.0k/s, took 3m40.513s, rate ~ 1.1M/s count: 252288000
2020-12-14T05:55:13.998Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 481.0k/s, took 4m46.082s, rate ~ 992.1k/s count: 283824000
2020-12-14T05:56:19.728Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 479.8k/s, took 5m51.812s, rate ~ 896.4k/s count: 315360000
2020-12-14T05:57:27.788Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 463.4k/s, took 6m59.872s, rate ~ 826.2k/s count: 346896000
2020-12-14T05:57:44.668Z - mysql(watt) <-> postgres took 7m16.752s, rate ~ 839.8k/s count: 366790994
2020-12-14T05:57:44.669Z - Verified mysql(watt) <-> postgres:
2020-12-14T05:57:44.669Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) Equal
2020-12-14T05:57:44.669Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-14T05:57:44.669Z - [2016-03-12T06:35:35Z, 2020-04-13T15:03:10Z](124627866) Equal
2020-12-14T05:57:44.669Z - [2020-04-13T15:03:11Z, 2020-11-20T23:32:33Z](18986075) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20200807.2218Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1919.074s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2020-08-07 22:18:41	134719305

- run mysql restore
2020-12-14T06:30:10.928Z - Starting TED1K mysql restore
2020-12-14T06:30:10.931Z - Connected to MySQL
2020-12-14T06:30:10.945Z - Connected to Postgres
2020-12-14T06:30:10.990Z - -=- mysql(watt) <-> postgres
2020-12-14T06:30:38.175Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.2M/s, took 27.184s, rate ~ 1.2M/s count: 31536000
2020-12-14T06:31:03.780Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.2M/s, took 52.79s, rate ~ 1.2M/s count: 63072000
2020-12-14T06:31:32.424Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.1M/s, took 1m21.434s, rate ~ 1.2M/s count: 94608000
2020-12-14T06:31:56.529Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.3M/s, took 1m45.539s, rate ~ 1.2M/s count: 126144000
2020-12-14T06:32:25.261Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.1M/s, took 2m14.27s, rate ~ 1.2M/s count: 157680000
2020-12-14T06:32:49.731Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.3M/s, took 2m38.741s, rate ~ 1.2M/s count: 189216000
2020-12-14T06:33:15.307Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 1.2M/s, took 3m4.317s, rate ~ 1.2M/s count: 220752000
2020-12-14T06:34:18.629Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 498.0k/s, took 4m7.639s, rate ~ 1.0M/s count: 252288000
2020-12-14T06:35:28.242Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 453.0k/s, took 5m17.252s, rate ~ 894.6k/s count: 283824000
2020-12-14T06:36:35.464Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 469.1k/s, took 6m24.474s, rate ~ 820.2k/s count: 315360000
2020-12-14T06:37:43.226Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 465.4k/s, took 7m32.236s, rate ~ 767.1k/s count: 346896000
2020-12-14T06:38:14.784Z - mysql(watt) <-> postgres took 8m3.793s, rate ~ 758.2k/s count: 366790994
2020-12-14T06:38:14.784Z - Verified mysql(watt) <-> postgres:
2020-12-14T06:38:14.784Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) Equal
2020-12-14T06:38:14.784Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-14T06:38:14.784Z - [2016-03-12T06:35:35Z, 2020-08-07T22:18:41Z](134644279) Equal
2020-12-14T06:38:14.784Z - [2020-08-07T22:18:42Z, 2020-11-20T23:32:33Z](8969662) MissingInA
- Done Restoring Database

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20201120.2332Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
ERROR 1114 (HY000) at line 3669: The table 'watt' is full
read unix @->/var/run/docker.sock: read: connection reset by peer
1798.592s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2020-06-01 20:30:47	128954558

- run mysql restore
# command-line-arguments
panic: no space left on device

goroutine 1 [running]:
cmd/link/internal/ld.Main(0x871840, 0x20, 0x20, 0x1, 0x7, 0x10, 0x0, 0x0, 0x6dbb5f, 0x1b, ...)
	/snap/go/current/src/cmd/link/internal/ld/main.go:320 +0x21bd
main.main()
	/snap/go/current/src/cmd/link/main.go:68 +0x1dc
- Done Restoring Database


real	1292m49.932s
user	226m43.214s
sys	60m46.526s

-=-= Restoring database from snapshot: ./data/archive/mirror/ted/ted.watt.20201120.2332Z.sql.bz2
- Drop tables watt and ted_native before restore, if present
- Show remaining tables, before restore
- Restoring database...
1950.910s

- Expect something recent in watt table
min(stamp)	max(stamp)	count(*)
2007-08-28 00:02:18	2020-11-20 23:32:33	143688967

- run mysql restore
2020-12-14T22:01:46.513Z - Starting TED1K mysql restore
2020-12-14T22:01:46.580Z - Connected to MySQL
2020-12-14T22:01:46.593Z - Connected to Postgres
2020-12-14T22:01:46.603Z - -=- mysql(watt) <-> postgres
2020-12-14T22:02:13.123Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.2M/s, took 26.52s, rate ~ 1.2M/s count: 31536000
2020-12-14T22:02:36.277Z - mysql(watt) <-> postgres (2010-10-15) inner rate ~ 1.4M/s, took 49.673s, rate ~ 1.3M/s count: 63072000
2020-12-14T22:03:05.671Z - mysql(watt) <-> postgres (2011-11-10) inner rate ~ 1.1M/s, took 1m19.068s, rate ~ 1.2M/s count: 94608000
2020-12-14T22:03:34.003Z - mysql(watt) <-> postgres (2012-11-14) inner rate ~ 1.1M/s, took 1m47.4s, rate ~ 1.2M/s count: 126144000
2020-12-14T22:04:04.259Z - mysql(watt) <-> postgres (2013-11-21) inner rate ~ 1.0M/s, took 2m17.656s, rate ~ 1.1M/s count: 157680000
2020-12-14T22:04:35.401Z - mysql(watt) <-> postgres (2014-12-14) inner rate ~ 1.0M/s, took 2m48.798s, rate ~ 1.1M/s count: 189216000
2020-12-14T22:05:08.462Z - mysql(watt) <-> postgres (2015-12-30) inner rate ~ 953.9k/s, took 3m21.859s, rate ~ 1.1M/s count: 220752000
2020-12-14T22:06:19.444Z - mysql(watt) <-> postgres (2017-02-19) inner rate ~ 444.3k/s, took 4m32.841s, rate ~ 924.7k/s count: 252288000
2020-12-14T22:07:36.439Z - mysql(watt) <-> postgres (2018-03-19) inner rate ~ 409.6k/s, took 5m49.836s, rate ~ 811.3k/s count: 283824000
2020-12-14T22:08:51.421Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 420.6k/s, took 7m4.817s, rate ~ 742.3k/s count: 315360000
2020-12-14T22:10:30.153Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 319.4k/s, took 8m43.55s, rate ~ 662.6k/s count: 346896000
2020-12-14T22:11:16.325Z - mysql(watt) <-> postgres took 9m29.722s, rate ~ 643.8k/s count: 366790994
2020-12-14T22:11:16.328Z - Verified mysql(watt) <-> postgres:
2020-12-14T22:11:16.328Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) Equal
2020-12-14T22:11:16.328Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-14T22:11:16.328Z - [2016-03-12T06:35:35Z, 2020-11-20T23:32:33Z](143613941) Equal
- Done Restoring Database


real	42m29.375s

So including the last running (1292m49.932s+42m29.375s) Phase 2 took 22 hours

```

## Phase 2 - Verification after cleanuo

After removing the 75026 rogue entries from 2007-08-28

```bash
$ time go run cmd/mysqlrestore/mysqlrestore.go 
2020-12-14T23:15:01.029Z - Starting TED1K mysql restore
2020-12-14T23:15:01.070Z - Connected to MySQL
2020-12-14T23:15:01.088Z - Connected to Postgres
2020-12-14T23:15:01.105Z - -=- mysql(watt) <-> postgres
2020-12-14T23:15:28.161Z - mysql(watt) <-> postgres (2009-08-21) inner rate ~ 1.2M/s, took 27.056s, rate ~ 1.2M/s count: 31536000
2020-12-14T23:15:49.876Z - mysql(watt) <-> postgres (2010-10-16) inner rate ~ 1.5M/s, took 48.771s, rate ~ 1.3M/s count: 63072000
2020-12-14T23:16:14.055Z - mysql(watt) <-> postgres (2011-11-11) inner rate ~ 1.3M/s, took 1m12.95s, rate ~ 1.3M/s count: 94608000
2020-12-14T23:16:41.839Z - mysql(watt) <-> postgres (2012-11-15) inner rate ~ 1.1M/s, took 1m40.734s, rate ~ 1.3M/s count: 126144000
2020-12-14T23:17:04.183Z - mysql(watt) <-> postgres (2013-11-22) inner rate ~ 1.4M/s, took 2m3.077s, rate ~ 1.3M/s count: 157680000
2020-12-14T23:17:29.451Z - mysql(watt) <-> postgres (2014-12-15) inner rate ~ 1.2M/s, took 2m28.345s, rate ~ 1.3M/s count: 189216000
2020-12-14T23:17:52.266Z - mysql(watt) <-> postgres (2015-12-31) inner rate ~ 1.4M/s, took 2m51.161s, rate ~ 1.3M/s count: 220752000
2020-12-14T23:18:58.806Z - mysql(watt) <-> postgres (2017-02-20) inner rate ~ 473.9k/s, took 3m57.701s, rate ~ 1.1M/s count: 252288000
2020-12-14T23:20:18.553Z - mysql(watt) <-> postgres (2018-03-20) inner rate ~ 395.4k/s, took 5m17.448s, rate ~ 894.1k/s count: 283824000
2020-12-14T23:21:37.802Z - mysql(watt) <-> postgres (2019-03-27) inner rate ~ 397.9k/s, took 6m36.697s, rate ~ 795.0k/s count: 315360000
2020-12-14T23:23:00.341Z - mysql(watt) <-> postgres (2020-04-03) inner rate ~ 382.1k/s, took 7m59.236s, rate ~ 723.9k/s count: 346896000
2020-12-14T23:23:54.516Z - mysql(watt) <-> postgres took 8m53.411s, rate ~ 687.5k/s count: 366715968
2020-12-14T23:23:54.521Z - Verified mysql(watt) <-> postgres:
2020-12-14T23:23:54.521Z - [2007-08-28T00:02:18Z, 2007-08-28T20:56:22Z](75026) MissingInB
2020-12-14T23:23:54.521Z - [2008-07-30T00:04:40Z, 2016-02-14T21:24:21Z](223102027) MissingInA
2020-12-14T23:23:54.522Z - [2016-03-12T06:35:35Z, 2020-11-20T23:32:33Z](143613941) Equal

real	8m57.163s
