package domain

import "errors"

var ErrEditTimeoutExceed = errors.New("уже нельзя внести изменения в комментарий")
var ErrEditNotMyComment = errors.New("можно изменять только свои комментарии")
var ErrCommentNotExist = errors.New("комментария с таким ID не существует")
