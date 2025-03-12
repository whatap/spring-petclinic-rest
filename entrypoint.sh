#!/bin/sh

WHATAP_JAR=$(ls /data/agent/whatap.agent-*.jar | awk -F'[.]|-' '{printf("%04d %04d %04d %04d %s\n", $(NF-4), $(NF-3), $(NF-2), $(NF-1), $0)}' | sort | tail -1 | awk '{print $NF}')
JAVA_OPTS="-javaagent:${WHATAP_JAR} -Dwhatap.okind=petclinic -Duser.timezone=GMT"
JAVA_OPTS="-javaagent:${WHATAP_JAR} -Dwhatap.oname=petclinic-01 -Duser.timezone=GMT"
JAVA_OPTS="${JAVA_OPTS} --add-opens=java.base/java.lang=ALL-UNNAMED"

if [ -n "$WHATAP_CONF" ]; then
        echo "$WHATAP_CONF" >/data/agent/whatap.conf
        echo "Saved WHATAP_CONF to /data/agent/whatap.conf"
else
        echo "WHATAP_CONF is not set. Skipping file write."
fi

java ${JAVA_OPTS} -jar /data/whatap/bin/app.jar
