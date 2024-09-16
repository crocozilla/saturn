PROG  START  COPY
      INTDEF PROG
      INTDEF UP
      INTDEF DOWN
LOOP  INTUSE
X     INTUSE
SIG   LOAD   UP
      CALL   LOOP
      ADD    1
      ADD    1,I
      ADD    #1
      BR     SIG
      COPY   1 2
UP    CONST  15
DOWN  CONST  SIG
      END