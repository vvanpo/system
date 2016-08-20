A formal language is specified as a set (often infinite) of all valid strings within said language.

Example - UTF-8
---------------

byte = {0,…,ff₁₆}

bytestream = {

𝕌 = {0,…,14fff₁₆,16000₁₆,…,18fff₁₆,1b000₁₆,…,1bfff₁₆,1d000₁₆,…,2cfff₁₆,2f000₁₆,…,2ffff₁₆,e0000₁₆,…,e0fff₁₆} # the set of all unicode code points

UTF-8 = {} # single UTF-8 character encoding
