PROG START SEC
     INTDEF PROG2
     INTDEF CAS
     INTDEF BRA
     INTDEF LOOP
UP   INTUSE
LOOP SUB CAS
     ADD UP
     READ LOC
     RET
BRA  CALL LOOP
     STORE X
     STOP
LOC  SPACE
CAS  CONST 4
X    SPACE
     END