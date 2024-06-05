#define ENABLE_WKE 1

#define __cdecl __cdecl 
#define __declspec __declspec 
#define __stdcall __stdcall 
#define __fastcall __fastcall

typedef unsigned short     wchar_t;
#define HAVE_WCHAR_T
#define _WCHAR_T_DEFINED


typedef long long int64;
#define __int64 int64

typedef wchar_t WCHAR;    // wc,   16-bit UNICODE character
typedef unsigned long       DWORD;
typedef unsigned short      WORD;
typedef void *HANDLE;
typedef WCHAR *LPWSTR;
typedef unsigned char       BYTE;
typedef BYTE             *LPBYTE;
typedef long LONG;
typedef DWORD   COLORREF;
#define DECLARE_HANDLE(name) struct name##__ { int unused; }; typedef struct name##__ *name
DECLARE_HANDLE(HINSTANCE);
DECLARE_HANDLE(HDC);
DECLARE_HANDLE(HWND);
typedef int                 BOOL;
typedef unsigned int        UINT;/* Types use for passing & returning polymorphic values */
typedef UINT WPARAM;
typedef LONG LPARAM;
typedef LONG LRESULT;


typedef struct _STARTUPINFOW {
    DWORD   cb;
    LPWSTR  lpReserved;
    LPWSTR  lpDesktop;
    LPWSTR  lpTitle;
    DWORD   dwX;
    DWORD   dwY;
    DWORD   dwXSize;
    DWORD   dwYSize;
    DWORD   dwXCountChars;
    DWORD   dwYCountChars;
    DWORD   dwFillAttribute;
    DWORD   dwFlags;
    WORD    wShowWindow;
    WORD    cbReserved2;
    LPBYTE  lpReserved2;
    HANDLE  hStdInput;
    HANDLE  hStdOutput;
    HANDLE  hStdError;
} STARTUPINFOW;

typedef struct tagRECT
{
    LONG    left;
    LONG    top;
    LONG    right;
    LONG    bottom;
} RECT;


typedef struct tagPOINT
{
    LONG  x;
    LONG  y;
} POINT;

