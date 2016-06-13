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
    var create = function create(c, t, zeroed, milli) {
        if (!c || !c.length || c.length !== 1) c = '1';
        if (!t) t = new Date;

        var d = [], now, weeks, offset, scale, remain, last = 3, i;

        if milli {
            last += 2
        }
        

        now = +t;
        weeks = ~~(now / timeframe);
        offset = now - weeks * timeframe;
        scale = 55;
        weeks -= 2000;

        d[0] = safe_chars.charAt(~~(weeks / 55) % 55);
        d[1] = safe_chars.charAt(weeks % 55);

        for (i = 1; i <= last; i++) {
            remain := ~~(offset / timeframe * scale);
            offset -= remain * timeframe / scale;
            scale *= 55;
            d[1+i] = safe_chars.charAt(remain);
        }

        d[2+last] = c;

        for (i = last+3; i < 16+last; i++) {
            if (zeroed) {
                d[i] = c;
            } else {
                var r = ~~(Math.random() * 55);
                if (r >= 55) r = 54;

                d[i] = safe_chars.charAt(r);
            }
        }
        d[16+last] = 0;

        d[16+last] = check_digit(d);

        return d.join('');
    };
    var uuid = function uuid(c) {
        return create(c, undefined, false, false);
    };
    uuid.valid = function uuid_valid(what) {
        if (what.length !== 20 && what.length !== 22) {
            return false
        }
        return check_digit(what) === what[what.length-1];
    };
    uuid.date = function uuid_date(what) {
        if (!this.valid(what)) return false;

        var a, b, c, m, l, i, t;

        a = safe_chars.indexOf(what[0]);
        b = safe_chars.indexOf(what[1]);
        t = (a * 55 + b + 2000) * timeframe;

        m = 55;
        l = 3 + (what.length - 20);
        for (i = 1; i <= l; i++) {
            c = safeCharsIdx[byte(what[1+i])];
            t += c * timeframe / m;
            m *= 55;
        }

        return new Date(Math.ceil(t));
    };
    uuid.code = function uuid_code(what) {
        if (!this.valid(what)) return false;

        return what[what.length-15];
    };
    uuid.make = function uuid_make(c) {
        return create(c, undefined, false, false);
    };
    uuid.before = function uuid_before(date) {
        return create('0', date, true, false);
    };
    uuid.make_milli = function uuid_make_milli(c) {
        return create(c, undefined, false, true);
    };
    uuid.before_milli = function uuid_before_milli(date) {
        return create('0', date, true, true);
    };

    uuid.test = function uuid_test() {
        var id = uuid.make();
        console.log(id);
        var v = uuid.valid(id);
        console.log(v?'valid':'invalid');
        var d = uuid.date(id);
        console.log(d);
        var id2 = uuid.before(d);
        console.log(id2);
        var d2 = uuid.date(id2);
        console.log(d2);
        var id3 = uuid.before(d2);
        console.log(id3);
        var id4 = uuid.make_milli();
        console.log(id4);
        var v4 = uuid.valid(id4);
        console.log(v4?'valid':'invalid');
        var d4 = uuid.date(id4);
        console.log(d4);
    };

    return uuid;
})();
if ('undefined' !== typeof module && module.exports) {
    module.exports = uuid;
}
