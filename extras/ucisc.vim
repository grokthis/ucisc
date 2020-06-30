" uCISC syntax file
" Language:    uCISC
" URL:         https://github.com/grokthis/ucisc-ruby
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

setlocal iskeyword+=-,?,<,>

syntax keyword uciscOpcode copy compute
highlight link uciscOpcode Function

syntax match uciscComment /#.*/
syntax match uciscComment /\/[^\/]*\//
syntax match uciscComment /'[^ ]*/
highlight link uciscComment Comment

syntax match uciscLabel /\([^a-zA-Z_:&@!]*\)\@=[a-zA-Z_:&@!][a-zA-Z0-9_:&@!]*:/
highlight link uciscLabel Identifier

syntax match uciscControl /[(){}]/
syntax match uciscControl /\(break\|loop\)/
highlight link uciscControl Statement

syntax match uciscImmediate /\<\((\)\?\(-\)\?\(0x\)\?[0-9a-fA-F]\+\([/.]\|\s\|,\)\@=/
highlight link uciscImmediate Number

syntax match uciscArg /\$[a-zA-Z0-9_:&@!?*]\+/
highlight link uciscArg Type

syntax keyword uciscOption push pop as
syntax match uciscOption /\(.\)\@=\<\(disp\|imm\)\>\(\/\|\s\|,\|)\|$\)\@=/
syntax match uciscOption /\(.\)\@=\<\(reg\|mem\|val\)\>\(\/\|\s\|,\|)\|$\)\@=/
syntax match uciscOption /\(.\)\@=\<\(op\|inc\|eff\)\>\(\/\|\s\|$\)\@=/
highlight link uciscOption Define

syntax match uciscData /^[ ]*% *\([0-9a-fA-F][0-9a-fA-F][ ]*\)*/
highlight link uciscData Number

set comments-=:#
set comments+=n:#
syntax match subxCommentedCode "#? .*"  | highlight link subxCommentedCode CommentedCode
let b:cmt_head = "#? "

let &cpo = s:save_cpo
