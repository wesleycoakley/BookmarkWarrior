#!/sbin/openrc-run

pidfile=/var/run/bmw.pid
name="Bookmark Warrior"
datadir=/usr/local/share/BookmarkWarrior
description="Bookmark Warrior"
runas="bookmarkwarrior"

depend() {
	need net
	need localmount
}

start() {
	ebegin "BookmarkWarrior is spinning up"
	start-stop-daemon --start --exec /usr/local/bin/BookmarkWarrior \
		--group "$runas" --user "$runas" --background --pidfile "$pidfile" \
		--chdir "$datadir"
	eend $?
}

stop() {
	ebegin "Stopping BookmarkWarrior"
	start-stop-daemon --stop --exec /usr/local/bin/BookmarkWarrior \
		--pidfile "$pidfile" --chdir "$datadir"
	eend $?
}
