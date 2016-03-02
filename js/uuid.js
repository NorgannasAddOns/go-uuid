/**
 * A compatible version for creating and decoding go-uuid compatible uuids in javascript.
 * Usage:
 *   var uuid = require('./uuid.js');
 */
var uuid = (function uuid(){
    var safe_chars = '23456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz';
    var timeframe = 604800000; // 1 week in ms
    var check_digit = function check_digit(what) {
        var hash = 0,
            n = what.length;
        for (var i = 0; i < n - 1; i++) {
            hash = ((hash << 5) - hash) + what[i].charCodeAt(0);
            hash |= 0;
        }
        hash = (hash >>> 0) % 55;
        if (hash >= 55) hash = 54;
        return safe_chars.charAt(hash);
    };
    var create = function create(c, t, zeroed) {
        if (!c || !c.length || c.length !== 1) c = '1';
        if (!t) t = new Date;

        var d = [], now, weeks, offset, scale, remain, remain2, remain3;

        now = +t;
        weeks = ~~(now / timeframe);
        offset = now - weeks * timeframe;
        scale = 55;
        remain = ~~(offset / timeframe * scale);
        offset -= remain * timeframe / scale;
        scale *= 55;
        remain2 = ~~(offset / timeframe * scale);
        offset -= remain2 * timeframe / scale;
        scale *= 55;
        remain3 = ~~(offset / timeframe * scale);

        weeks -= 2000;

        d[0] = safe_chars.charAt(~~(weeks / 55) % 55);
        d[1] = safe_chars.charAt(weeks % 55);
        d[2] = safe_chars.charAt(remain);
        d[3] = safe_chars.charAt(remain2);
        d[4] = safe_chars.charAt(remain3);
        d[5] = c;

        for (var i = 6; i < 19; i++) {
            if (zeroed) {
                d[i] = c;
            } else {
                var r = ~~(Math.random() * 55);
                if (r >= 55) r = 54;

                d[i] = safe_chars.charAt(r);
            }
        }
        d[19] = 0;

        d[19] = check_digit(d);

        return d.join('');
    };
    var uuid = function uuid(c) {
        return create(c);
    };
    uuid.valid = function uuid_valid(what) {
        return check_digit(what) === what[what.length-1];
    };
    uuid.date = function uuid_date(what) {
        if (what.length !== 20 || !this.valid(what)) return false;

        var a, b, c, d, e, t;
        a = safe_chars.indexOf(what[0]);
        b = safe_chars.indexOf(what[1]);
        c = safe_chars.indexOf(what[2]);
        d = safe_chars.indexOf(what[3]);
        e = safe_chars.indexOf(what[4]);
        t = (a * 55 + b + 2000) * timeframe;
        t += (c * timeframe / 55);
        t += (d * timeframe / (55 * 55));
        t += (e * timeframe / (55 * 55 * 55));

        return new Date(Math.ceil(t));
    };
    uuid.code = function uuid_code(what) {
        if (what.length !== 20 || !this.valid(what)) return false;

        return what[5];
    };
    uuid.make = function uuid_make(c) {
        return create(c);
    };
    uuid.before = function uuid_before(date) {
        return create('0', date, true);
    };

    uuid.test = function uuid_test() {
        id = uuid.make();
        console.log(id);
        v = uuid.valid(id);
        console.log(v?'valid':'invalid');
        d = uuid.date(id);
        console.log(d);
        id2 = uuid.before(d);
        console.log(id2);
        d2 = uuid.date(id2);
        console.log(d2);
        id3 = uuid.before(d2);
        console.log(id3);
    };

    return uuid;
})();
if ('undefined' !== typeof module && module.exports) {
    module.exports = uuid;
}
