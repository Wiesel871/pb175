AC_INIT([pb175], [1.0], [536388@mail.muni.cz])
AM_INIT_AUTOMAKE

AC_ARG_WITH([name],
            [--with-name],
            [AC_SUBST([name], "$withval")])
AC_ARG_WITH([email],
            [--with-email],
            [AC_SUBST([email], "$withval")])
AC_ARG_WITH([password],
            [--with-password],
            [AC_SUBST([password], "$withval")])

AC_ARG_WITH([prefix],
            [--prefix],
            [PREFIX="$withval"],
            [prefix="."])

AC_SUBST([PREFIX])
AC_CONFIG_FILES([main], [chmod +x main])
AC_CONFIG_FILES([Makefile])

AC_OUTPUT
