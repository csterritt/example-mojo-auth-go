#!/usr/bin/env perl
$l1 = <>;
$l2 = <>;
$l3 = <>;
while (<>) {
    $a = $_;
    if ($l1 =~ /^\s*\],/ && $l2 =~ /tls_connection_policies/) {
        print("]\n");
        $l1 = <>;
        $l2 = <>;
        $l3 = <>;
        continue;
    }

    print($l1);
    $l1 = $l2;
    $l2 = $l3;
    $l3 = $a;
}
print($l1);
print($l2);
print($l3);
