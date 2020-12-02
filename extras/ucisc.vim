" uCISC syntax file
" Language:    uCISC
" URL:         https://github.com/grokthis/ucisc
" License:     MIT
"
" Copy this into your ftplugin directory, and add the following to your vimrc
" or to .vim/ftdetect/ucisc.vim:
"   autocmd BufReadPost,BufNewFile *.ucisc set filetype=ucisc

let s:save_cpo = &cpo
set cpo&vim

" setlocal iskeyword=@,48-57,?,!,_,$,-
setlocal formatoptions-=t  " allow long lines
setlocal formatoptions+=c  " but comments should still wrap

setlocal iskeyword+=-,?,<,>,!,~

syntax keyword uciscOpcode <- <~ <0 <0? <! <!? <p <p? <n <n? <o <o? <i <i?
highlight link uciscOpcode Type

syntax match uciscComment /#.*/
highlight link uciscComment Comment

syntax keyword uciscLabel copy and or xor inv shl shr swap lsb msb
syntax keyword uciscLabel add sub mult multu addc
highlight link uciscLabel Identifier

syntax keyword uciscControl pc
syntax match uciscControl /[{}]/
syntax match uciscControl /\(break\|loop\)/
highlight link uciscControl Statement

syntax keyword uciscArg push pop ->
highlight link uciscArg Function

syntax keyword uciscOption def var fun
highlight link uciscOption Define

syntax match uciscData /^[ ]*% *\([0-9a-fA-F][0-9a-fA-F][ ]*\)*/
highlight link uciscData Number

set comments-=:#
set comments+=n:#
syntax match subxCommentedCode "#? .*"  | highlight link subxCommentedCode CommentedCode
let b:cmt_head = "#? "

let &cpo = s:save_cpo
