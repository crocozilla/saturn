PROG  START  FIRST
      INTDEF PROG
      INTDEF UP
      INTDEF DOWN
LOOP  INTUSE
X     INTUSE
SIG   LOAD   UP
      CALL   LOOP
      BR     SIG
UP    CONST  15
DOWN  CONST  SIG
      END