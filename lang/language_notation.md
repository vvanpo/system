A formal language is specified as a set (often infinite) of all valid strings within said language.

Example - UTF-8
---------------

𝔹 = {0,…,ff₁₆} # set of bytes

𝕊 = {(b₁,…,bₙ) ∈ 𝔹ⁿ | n ∈ ℕ} # set of bytestreams

𝕌 = {0,…,14fff₁₆,16000₁₆,…,18fff₁₆,1b000₁₆,…,1bfff₁₆,1d000₁₆,…,2cfff₁₆,2f000₁₆,…,2ffff₁₆,e0000₁₆,…,e0fff₁₆} # the set of all valid unicode code points

f:  𝕌 → 𝕊
    u ↦ ({u | u ≤ 7f₁₆})

UTF-8 = {u ∈ 𝕌 | f(u)} # single UTF-8 character encoding
