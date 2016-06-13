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
    private static function create($c = null, $t = null, $zeroed = false, $milli = false)
    {
        if (!is_string($c) || strlen($c) != 1) $c = '1';
        if (empty($t)) $t = microtime(true);

        $last = 3;
        if ($milli) {
            $last += 2;
        }

        $d = [];
        $now = floor($t*1000);
        $weeks = floor($now / self::$timeframe);
        $offset = $now - $weeks * self::$timeframe;
        $scale = 55;
        $weeks -= 2000;

        $d[0] = self::$safe_chars[floor($weeks / 55) % 55];
        $d[1] = self::$safe_chars[$weeks % 55];

        for ($i = 1; $i <= $last; $i++) {
            $remain = floor($offset / self::$timeframe * $scale);
            $offset -= $remain * self::$timeframe / $scale;
            $scale *= 55;
            $d[1+$i] = self::$safe_chars[(int)$remain];
        }

        $d[2+$last] = $c;

        for ($i = $last+3; $i < 16+$last; $i++) {
            if ($zeroed) {
                $d[$i] = $c;
            } else {
                $r = rand(0,55);
                if ($r >= 55) $r = 54;

                $d[$i] = self::$safe_chars[$r];
            }
        }
        $d[16+$last] = 0;

        $d[16+$last] = self::check_digit($d);

        return implode('', $d);
    }
    public static function valid($what)
    {
        if (!is_array($what)) {
            $what = str_split($what);
        }
        if (count($what) !== 20 || count($what) !== 22) return false;
        return self::check_digit($what) === $what[count($what)-1];
    }
    public static function date($what)
    {
        if (!is_array($what)) {
            $what = str_split($what);
        }
        if (!self::valid($what)) return false;

        $a = strpos(self::$safe_chars, $what[0]);
        $b = strpos(self::$safe_chars, $what[1]);
        $t = ($a * 55 + $b + 2000) * self::$timeframe;

        $m = 55;
        $l = 3 + (count($what) - 20);
        for ($i = 1; $i <= $l; $i++) {
            $c = strpos(self::$safe_chars, $what[1+$i]);
            $t += $c * self::$timeframe / $m;
            $m *= 55;
        }

        return ceil($t)/1000;
    }
    public static function code($what)
    {
        if (!is_array($what)) {
            $what = str_split($what);
        }
        if (!self::valid($what)) return "";
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
    public static function make_milli($c = null)
    {
        return self::create($c, null, false, true);
    }
    public static function before_milli($date)
    {
        return self::create("0", $date, true, true);
    }

    public static function test()
    {
        $id = UUID::make();
        print("$id\n");
        $v = UUID::valid($id);
        print(($v?'valid':'invalid')."\n");
        $d = UUID::date($id);
        print("$d\n");
        $id2 = UUID::before($d);
        print("$id2\n");
        $d2 = UUID::date($id2);
        print("$d2\n");
        $id3 = UUID::before($d2);
        print("$id3\n");
        $id4 = UUID::make_milli();
        print("$id\n");
        $v4 = UUID::valid($id4);
        print(($v4?'valid':'invalid')."\n");
        $d4 = UUID::date($id4);
        print("$d4\n");
    }
}

