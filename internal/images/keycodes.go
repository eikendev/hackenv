package images

//revive:disable

const (
	KEY_RESERVED         uint = 0x0
	KEY_ESC              uint = 0x1
	KEY_1                uint = 0x2
	KEY_2                uint = 0x3
	KEY_3                uint = 0x4
	KEY_4                uint = 0x5
	KEY_5                uint = 0x6
	KEY_6                uint = 0x7
	KEY_7                uint = 0x8
	KEY_8                uint = 0x9
	KEY_9                uint = 0xa
	KEY_0                uint = 0xb
	KEY_MINUS            uint = 0xc
	KEY_EQUAL            uint = 0xd
	KEY_BACKSPACE        uint = 0xe
	KEY_TAB              uint = 0xf
	KEY_Q                uint = 0x10
	KEY_W                uint = 0x11
	KEY_E                uint = 0x12
	KEY_R                uint = 0x13
	KEY_T                uint = 0x14
	KEY_Y                uint = 0x15
	KEY_U                uint = 0x16
	KEY_I                uint = 0x17
	KEY_O                uint = 0x18
	KEY_P                uint = 0x19
	KEY_LEFTBRACE        uint = 0x1a
	KEY_RIGHTBRACE       uint = 0x1b
	KEY_ENTER            uint = 0x1c
	KEY_LEFTCTRL         uint = 0x1d
	KEY_A                uint = 0x1e
	KEY_S                uint = 0x1f
	KEY_D                uint = 0x20
	KEY_F                uint = 0x21
	KEY_G                uint = 0x22
	KEY_H                uint = 0x23
	KEY_J                uint = 0x24
	KEY_K                uint = 0x25
	KEY_L                uint = 0x26
	KEY_SEMICOLON        uint = 0x27
	KEY_APOSTROPHE       uint = 0x28
	KEY_GRAVE            uint = 0x29
	KEY_LEFTSHIFT        uint = 0x2a
	KEY_BACKSLASH        uint = 0x2b
	KEY_Z                uint = 0x2c
	KEY_X                uint = 0x2d
	KEY_C                uint = 0x2e
	KEY_V                uint = 0x2f
	KEY_B                uint = 0x30
	KEY_N                uint = 0x31
	KEY_M                uint = 0x32
	KEY_COMMA            uint = 0x33
	KEY_DOT              uint = 0x34
	KEY_SLASH            uint = 0x35
	KEY_RIGHTSHIFT       uint = 0x36
	KEY_KPASTERISK       uint = 0x37
	KEY_LEFTALT          uint = 0x38
	KEY_SPACE            uint = 0x39
	KEY_CAPSLOCK         uint = 0x3a
	KEY_F1               uint = 0x3b
	KEY_F2               uint = 0x3c
	KEY_F3               uint = 0x3d
	KEY_F4               uint = 0x3e
	KEY_F5               uint = 0x3f
	KEY_F6               uint = 0x40
	KEY_F7               uint = 0x41
	KEY_F8               uint = 0x42
	KEY_F9               uint = 0x43
	KEY_F10              uint = 0x44
	KEY_NUMLOCK          uint = 0x45
	KEY_SCROLLLOCK       uint = 0x46
	KEY_KP7              uint = 0x47
	KEY_KP8              uint = 0x48
	KEY_KP9              uint = 0x49
	KEY_KPMINUS          uint = 0x4a
	KEY_KP4              uint = 0x4b
	KEY_KP5              uint = 0x4c
	KEY_KP6              uint = 0x4d
	KEY_KPPLUS           uint = 0x4e
	KEY_KP1              uint = 0x4f
	KEY_KP2              uint = 0x50
	KEY_KP3              uint = 0x51
	KEY_KP0              uint = 0x52
	KEY_KPDOT            uint = 0x53
	KEY_ZENKAKUHANKAKU   uint = 0x55
	KEY_102ND            uint = 0x56
	KEY_F11              uint = 0x57
	KEY_F12              uint = 0x58
	KEY_RO               uint = 0x59
	KEY_KATAKANA         uint = 0x5a
	KEY_HIRAGANA         uint = 0x5b
	KEY_HENKAN           uint = 0x5c
	KEY_KATAKANAHIRAGANA uint = 0x5d
	KEY_MUHENKAN         uint = 0x5e
	KEY_KPJPCOMMA        uint = 0x5f
	KEY_KPENTER          uint = 0x60
	KEY_RIGHTCTRL        uint = 0x61
	KEY_KPSLASH          uint = 0x62
	KEY_SYSRQ            uint = 0x63
	KEY_RIGHTALT         uint = 0x64
	KEY_LINEFEED         uint = 0x65
	KEY_HOME             uint = 0x66
	KEY_UP               uint = 0x67
	KEY_PAGEUP           uint = 0x68
	KEY_LEFT             uint = 0x69
	KEY_RIGHT            uint = 0x6a
	KEY_END              uint = 0x6b
	KEY_DOWN             uint = 0x6c
	KEY_PAGEDOWN         uint = 0x6d
	KEY_INSERT           uint = 0x6e
	KEY_DELETE           uint = 0x6f
	KEY_MACRO            uint = 0x70
	KEY_MUTE             uint = 0x71
	KEY_VOLUMEDOWN       uint = 0x72
	KEY_VOLUMEUP         uint = 0x73
	KEY_POWER            uint = 0x74
	KEY_KPEQUAL          uint = 0x75
	KEY_KPPLUSMINUS      uint = 0x76
	KEY_PAUSE            uint = 0x77
	KEY_SCALE            uint = 0x78
	KEY_KPCOMMA          uint = 0x79
	KEY_HANGEUL          uint = 0x7a
	KEY_HANJA            uint = 0x7b
	KEY_YEN              uint = 0x7c
	KEY_LEFTMETA         uint = 0x7d
	KEY_RIGHTMETA        uint = 0x7e
	KEY_COMPOSE          uint = 0x7f
	KEY_STOP             uint = 0x80
	KEY_AGAIN            uint = 0x81
	KEY_PROPS            uint = 0x82
	KEY_UNDO             uint = 0x83
	KEY_FRONT            uint = 0x84
	KEY_COPY             uint = 0x85
	KEY_OPEN             uint = 0x86
	KEY_PASTE            uint = 0x87
	KEY_FIND             uint = 0x88
	KEY_CUT              uint = 0x89
	KEY_HELP             uint = 0x8a
	KEY_MENU             uint = 0x8b
	KEY_CALC             uint = 0x8c
	KEY_SETUP            uint = 0x8d
	KEY_SLEEP            uint = 0x8e
	KEY_WAKEUP           uint = 0x8f
	KEY_FILE             uint = 0x90
	KEY_SENDFILE         uint = 0x91
	KEY_DELETEFILE       uint = 0x92
	KEY_XFER             uint = 0x93
	KEY_PROG1            uint = 0x94
	KEY_PROG2            uint = 0x95
	KEY_WWW              uint = 0x96
	KEY_MSDOS            uint = 0x97
	KEY_SCREENLOCK       uint = 0x98
	KEY_DIRECTION        uint = 0x99
	KEY_CYCLEWINDOWS     uint = 0x9a
	KEY_MAIL             uint = 0x9b
	KEY_BOOKMARKS        uint = 0x9c
	KEY_COMPUTER         uint = 0x9d
	KEY_BACK             uint = 0x9e
	KEY_FORWARD          uint = 0x9f
	KEY_CLOSECD          uint = 0xa0
	KEY_EJECTCD          uint = 0xa1
	KEY_EJECTCLOSECD     uint = 0xa2
	KEY_NEXTSONG         uint = 0xa3
	KEY_PLAYPAUSE        uint = 0xa4
	KEY_PREVIOUSSONG     uint = 0xa5
	KEY_STOPCD           uint = 0xa6
	KEY_RECORD           uint = 0xa7
	KEY_REWIND           uint = 0xa8
	KEY_PHONE            uint = 0xa9
	KEY_ISO              uint = 0xaa
	KEY_CONFIG           uint = 0xab
	KEY_HOMEPAGE         uint = 0xac
	KEY_REFRESH          uint = 0xad
	KEY_EXIT             uint = 0xae
	KEY_MOVE             uint = 0xaf
	KEY_EDIT             uint = 0xb0
	KEY_SCROLLUP         uint = 0xb1
	KEY_SCROLLDOWN       uint = 0xb2
	KEY_KPLEFTPAREN      uint = 0xb3
	KEY_KPRIGHTPAREN     uint = 0xb4
	KEY_NEW              uint = 0xb5
	KEY_REDO             uint = 0xb6
	KEY_F13              uint = 0xb7
	KEY_F14              uint = 0xb8
	KEY_F15              uint = 0xb9
	KEY_F16              uint = 0xba
	KEY_F17              uint = 0xbb
	KEY_F18              uint = 0xbc
	KEY_F19              uint = 0xbd
	KEY_F20              uint = 0xbe
	KEY_F21              uint = 0xbf
	KEY_F22              uint = 0xc0
	KEY_F23              uint = 0xc1
	KEY_F24              uint = 0xc2
	KEY_PLAYCD           uint = 0xc8
	KEY_PAUSECD          uint = 0xc9
	KEY_PROG3            uint = 0xca
	KEY_PROG4            uint = 0xcb
	KEY_DASHBOARD        uint = 0xcc
	KEY_SUSPEND          uint = 0xcd
	KEY_CLOSE            uint = 0xce
	KEY_PLAY             uint = 0xcf
	KEY_FASTFORWARD      uint = 0xd0
	KEY_BASSBOOST        uint = 0xd1
	KEY_PRINT            uint = 0xd2
	KEY_HP               uint = 0xd3
	KEY_CAMERA           uint = 0xd4
	KEY_SOUND            uint = 0xd5
	KEY_QUESTION         uint = 0xd6
	KEY_EMAIL            uint = 0xd7
	KEY_CHAT             uint = 0xd8
	KEY_SEARCH           uint = 0xd9
	KEY_CONNECT          uint = 0xda
	KEY_FINANCE          uint = 0xdb
	KEY_SPORT            uint = 0xdc
	KEY_SHOP             uint = 0xdd
	KEY_ALTERASE         uint = 0xde
	KEY_CANCEL           uint = 0xdf
	KEY_BRIGHTNESSDOWN   uint = 0xe0
	KEY_BRIGHTNESSUP     uint = 0xe1
	KEY_MEDIA            uint = 0xe2
	KEY_SWITCHVIDEOMODE  uint = 0xe3
	KEY_KBDILLUMTOGGLE   uint = 0xe4
	KEY_KBDILLUMDOWN     uint = 0xe5
	KEY_KBDILLUMUP       uint = 0xe6
	KEY_SEND             uint = 0xe7
	KEY_REPLY            uint = 0xe8
	KEY_FORWARDMAIL      uint = 0xe9
	KEY_SAVE             uint = 0xea
	KEY_DOCUMENTS        uint = 0xeb
	KEY_BATTERY          uint = 0xec
	KEY_BLUETOOTH        uint = 0xed
	KEY_WLAN             uint = 0xee
	KEY_UWB              uint = 0xef
	KEY_UNKNOWN          uint = 0xf0
	KEY_VIDEO_NEXT       uint = 0xf1
	KEY_VIDEO_PREV       uint = 0xf2
	KEY_BRIGHTNESS_CYCLE uint = 0xf3
	KEY_BRIGHTNESS_ZERO  uint = 0xf4
	KEY_DISPLAY_OFF      uint = 0xf5
	KEY_WIMAX            uint = 0xf6
	KEY_OK               uint = 0x160
	KEY_SELECT           uint = 0x161
	KEY_GOTO             uint = 0x162
	KEY_CLEAR            uint = 0x163
	KEY_POWER2           uint = 0x164
	KEY_OPTION           uint = 0x165
	KEY_INFO             uint = 0x166
	KEY_TIME             uint = 0x167
	KEY_VENDOR           uint = 0x168
	KEY_ARCHIVE          uint = 0x169
	KEY_PROGRAM          uint = 0x16a
	KEY_CHANNEL          uint = 0x16b
	KEY_FAVORITES        uint = 0x16c
	KEY_EPG              uint = 0x16d
	KEY_PVR              uint = 0x16e
	KEY_MHP              uint = 0x16f
	KEY_LANGUAGE         uint = 0x170
	KEY_TITLE            uint = 0x171
	KEY_SUBTITLE         uint = 0x172
	KEY_ANGLE            uint = 0x173
	KEY_ZOOM             uint = 0x174
	KEY_MODE             uint = 0x175
	KEY_KEYBOARD         uint = 0x176
	KEY_SCREEN           uint = 0x177
	KEY_PC               uint = 0x178
	KEY_TV               uint = 0x179
	KEY_TV2              uint = 0x17a
	KEY_VCR              uint = 0x17b
	KEY_VCR2             uint = 0x17c
	KEY_SAT              uint = 0x17d
	KEY_SAT2             uint = 0x17e
	KEY_CD               uint = 0x17f
	KEY_TAPE             uint = 0x180
	KEY_RADIO            uint = 0x181
	KEY_TUNER            uint = 0x182
	KEY_PLAYER           uint = 0x183
	KEY_TEXT             uint = 0x184
	KEY_DVD              uint = 0x185
	KEY_AUX              uint = 0x186
	KEY_MP3              uint = 0x187
	KEY_AUDIO            uint = 0x188
	KEY_VIDEO            uint = 0x189
	KEY_DIRECTORY        uint = 0x18a
	KEY_LIST             uint = 0x18b
	KEY_MEMO             uint = 0x18c
	KEY_CALENDAR         uint = 0x18d
	KEY_RED              uint = 0x18e
	KEY_GREEN            uint = 0x18f
	KEY_YELLOW           uint = 0x190
	KEY_BLUE             uint = 0x191
	KEY_CHANNELUP        uint = 0x192
	KEY_CHANNELDOWN      uint = 0x193
	KEY_FIRST            uint = 0x194
	KEY_LAST             uint = 0x195
	KEY_AB               uint = 0x196
	KEY_NEXT             uint = 0x197
	KEY_RESTART          uint = 0x198
	KEY_SLOW             uint = 0x199
	KEY_SHUFFLE          uint = 0x19a
	KEY_BREAK            uint = 0x19b
	KEY_PREVIOUS         uint = 0x19c
	KEY_DIGITS           uint = 0x19d
	KEY_TEEN             uint = 0x19e
	KEY_TWEN             uint = 0x19f
	KEY_VIDEOPHONE       uint = 0x1a0
	KEY_GAMES            uint = 0x1a1
	KEY_ZOOMIN           uint = 0x1a2
	KEY_ZOOMOUT          uint = 0x1a3
	KEY_ZOOMRESET        uint = 0x1a4
	KEY_WORDPROCESSOR    uint = 0x1a5
	KEY_EDITOR           uint = 0x1a6
	KEY_SPREADSHEET      uint = 0x1a7
	KEY_GRAPHICSEDITOR   uint = 0x1a8
	KEY_PRESENTATION     uint = 0x1a9
	KEY_DATABASE         uint = 0x1aa
	KEY_NEWS             uint = 0x1ab
	KEY_VOICEMAIL        uint = 0x1ac
	KEY_ADDRESSBOOK      uint = 0x1ad
	KEY_MESSENGER        uint = 0x1ae
	KEY_DISPLAYTOGGLE    uint = 0x1af
	KEY_SPELLCHECK       uint = 0x1b0
	KEY_LOGOFF           uint = 0x1b1
	KEY_DOLLAR           uint = 0x1b2
	KEY_EURO             uint = 0x1b3
	KEY_FRAMEBACK        uint = 0x1b4
	KEY_FRAMEFORWARD     uint = 0x1b5
	KEY_CONTEXT_MENU     uint = 0x1b6
	KEY_MEDIA_REPEAT     uint = 0x1b7
	KEY_DEL_EOL          uint = 0x1c0
	KEY_DEL_EOS          uint = 0x1c1
	KEY_INS_LINE         uint = 0x1c2
	KEY_DEL_LINE         uint = 0x1c3
	KEY_FN               uint = 0x1d0
	KEY_FN_ESC           uint = 0x1d1
	KEY_FN_F1            uint = 0x1d2
	KEY_FN_F2            uint = 0x1d3
	KEY_FN_F3            uint = 0x1d4
	KEY_FN_F4            uint = 0x1d5
	KEY_FN_F5            uint = 0x1d6
	KEY_FN_F6            uint = 0x1d7
	KEY_FN_F7            uint = 0x1d8
	KEY_FN_F8            uint = 0x1d9
	KEY_FN_F9            uint = 0x1da
	KEY_FN_F10           uint = 0x1db
	KEY_FN_F11           uint = 0x1dc
	KEY_FN_F12           uint = 0x1dd
	KEY_FN_1             uint = 0x1de
	KEY_FN_2             uint = 0x1df
	KEY_FN_D             uint = 0x1e0
	KEY_FN_E             uint = 0x1e1
	KEY_FN_F             uint = 0x1e2
	KEY_FN_S             uint = 0x1e3
	KEY_FN_B             uint = 0x1e4
	KEY_BRL_DOT1         uint = 0x1f1
	KEY_BRL_DOT2         uint = 0x1f2
	KEY_BRL_DOT3         uint = 0x1f3
	KEY_BRL_DOT4         uint = 0x1f4
	KEY_BRL_DOT5         uint = 0x1f5
	KEY_BRL_DOT6         uint = 0x1f6
	KEY_BRL_DOT7         uint = 0x1f7
	KEY_BRL_DOT8         uint = 0x1f8
	KEY_BRL_DOT9         uint = 0x1f9
	KEY_BRL_DOT10        uint = 0x1fa
	KEY_NUMERIC_0        uint = 0x200
	KEY_NUMERIC_1        uint = 0x201
	KEY_NUMERIC_2        uint = 0x202
	KEY_NUMERIC_3        uint = 0x203
	KEY_NUMERIC_4        uint = 0x204
	KEY_NUMERIC_5        uint = 0x205
	KEY_NUMERIC_6        uint = 0x206
	KEY_NUMERIC_7        uint = 0x207
	KEY_NUMERIC_8        uint = 0x208
	KEY_NUMERIC_9        uint = 0x209
	KEY_NUMERIC_STAR     uint = 0x20a
	KEY_NUMERIC_POUND    uint = 0x20b
	KEY_RFKILL           uint = 0x20c
)

//revive:enable
