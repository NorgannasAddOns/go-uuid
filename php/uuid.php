<?php
/**
 * A compatible version for creating and decoding go-uuid compatible uuid's in php.
 */
class UUID
{
    private static $safe_chars = '23456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz';
    private static $timeframe = 604800000; // 1 week in ms
    private static function bitwise32($num, $zeroFill = false)
    {
        $bin = decbin($num);
        $bin = str_pad($bin, 32, "0", STR_PAD_LEFT);
        $bin = substr($bin, strlen($bin) - 32);
        if ($bin{0} == "1" && !$zeroFill) {
            return -(pow(2, 31) - bindec(substr($bin, 1)));
        }
        return bindec($bin);
    }
    private static function check_digit($what)
    {
        $hash = 0;
        $n = count($what);
        for ($i = 0; $i < $n - 1; $i++) {
            $hash = (self::bitwise32($hash<<5) - $hash) + ord($what[$i]);
            $hash = self::bitwise32($hash | 0);
        }
        $hash = self::bitwise32($hash>>0, true) % 55;
        if ($hash >= 55) $hash = 54;
        return self::$safe_chars[$hash];
    }
    private static function create($c = null, $t = null, $zeroed = false)
    {
        if (!is_string($c) || strlen($c) != 1) $c = '1';
        if (empty($t)) $t = microtime(true);

        $d = [];
        $now = floor($t*1000);
        $weeks = floor($now / self::$timeframe);
        $offset = $now - $weeks * self::$timeframe;
        $scale = 55;
        $remain = floor($offset / self::$timeframe * $scale);
        $offset -= $remain * self::$timeframe / $scale;
        $scale *= 55;
        $remain2 = floor($offset / self::$timeframe * $scale);
        $offset -= $remain2 * self::$timeframe / $scale;
        $scale *= 55;
        $remain3 = floor($offset / self::$timeframe * $scale);

        $weeks -= 2000;

        $d[0] = self::$safe_chars[floor($weeks / 55) % 55];
        $d[1] = self::$safe_chars[$weeks % 55];
        $d[2] = self::$safe_chars[(int)$remain];
        $d[3] = self::$safe_chars[(int)$remain2];
        $d[4] = self::$safe_chars[(int)$remain3];
        $d[5] = $c;

        for ($i = 6; $i < 19; $i++) {
            if ($zeroed) {
                $d[$i] = $c;
            } else {
                $r = rand(0,55);
                if ($r >= 55) $r = 54;

                $d[$i] = self::$safe_chars[$r];
            }
        }
        $d[19] = 0;

        $d[19] = self::check_digit($d);

        return implode('', $d);
    }
    public static function valid($what)
    {
        if (!is_array($what)) {
            $what = str_split($what);
        }
        return self::check_digit($what) === $what[count($what)-1];
    }
    public static function date($what)
    {
        if (!is_array($what)) {
            $what = str_split($what);
        }
        if (count($what) !== 20 || !self::valid($what)) return false;

        $a = strpos(self::$safe_chars, $what[0]);
        $b = strpos(self::$safe_chars, $what[1]);
        $c = strpos(self::$safe_chars, $what[2]);
        $d = strpos(self::$safe_chars, $what[3]);
        $e = strpos(self::$safe_chars, $what[4]);
        $t = ($a * 55 + $b + 2000) * self::$timeframe;
        $t += ($c * self::$timeframe / 55);
        $t += ($d * self::$timeframe / (55 * 55));
        $t += ($e * self::$timeframe / (55 * 55 * 55));

        return ceil($t)/1000;
    }
    public static function code($what)
    {
        if (!is_array($what)) {
            $what = str_split($what);
        }
        if (count($what) !== 20 || !self::valid($what)) return "";
        return $what[5];
    }
    public static function make($c = null)
    {
        return self::create($c);
    }
    public static function before($date)
    {
        return self::create("0", $date, true);
    }

    public static function test()
    {
        $id = MongoishUUID::make();
        print("$id\n");
        $v = MongoishUUID::valid($id);
        print(($v?'valid':'invalid')."\n");
        $d = MongoishUUID::date($id);
        print("$d\n");
        $id2 = MongoishUUID::before($d);
        print("$id2\n");
        $d2 = MongoishUUID::date($id2);
        print("$d2\n");
        $id3 = MongoishUUID::before($d2);
        print("$id3\n");
    }
}

