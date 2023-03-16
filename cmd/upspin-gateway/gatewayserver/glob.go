package gatewayserver

import (
	"fmt"
	"strings"

	"upspin.io/path"
	"upspin.io/upspin"
	"upspin.io/user"
)

// HasGlobChar reports whether the string contains an unescaped Glob metacharacter.
func HasGlobChar(pattern string) bool {
	esc := false
	for _, r := range pattern {
		if esc {
			esc = false // TODO: What if next rune is '/'?
			continue
		}
		switch r {
		case '\\':
			esc = true
		case '*', '[', '?':
			return true
		}
	}
	return false
}

// AtSign processes a leading at sign, if any, in the Upspin file name and replaces it
// with the current user name. The name must be strictly "@" or begin with "@/",
// possibly with a suffix ("@+snapshot", "@+camera/").
// (If the user name ready has a suffix and the file starts @+suffix,
// the returned user name appends the suffix from the argument.)
// Unlike Tilde, it does not look up others user's roots.
// The argument is of type string; once a file becomes an upspin.PathName it should
// not be passed to this function.
// If the file name does not begin with an at sign, AtSign returns the argument
// unchanged except for promotion to upspin.PathName.
// If the target user does not exist, it returns the original string.
func (s *Gateway) AtSign(file string) upspin.PathName {
	if s.cfg == nil || file == "" || file[0] != '@' {
		return upspin.PathName(file)
	}
	userStr := string(s.cfg.UserName())
	if file == "@" {
		return upspin.PathName(userStr + "/")
	}
	if strings.HasPrefix(file, "@/") {
		return upspin.PathName(userStr + file[1:])
	}
	if strings.HasPrefix(file, "@+") {
		// Need to split the user name.
		usr, _, domain, err := user.Parse(s.cfg.UserName())
		if err != nil { // Can't happen.
			return upspin.PathName(file)
		}
		slash := strings.IndexByte(file, '/')
		if slash < 0 {
			return upspin.PathName(usr + file[1:] + "@" + domain + "/")
		}
		return upspin.PathName(usr + file[1:slash] + "@" + domain + file[slash:])
	}
	return upspin.PathName(file)
}

// GlobUpspinPath glob-expands the argument, which must be a syntactically
// valid Upspin glob pattern (including a plain path name). It returns just
// the path names. If the pattern matches no paths, the function exits.
func (s *Gateway) GlobUpspinPath(pattern string) ([]upspin.PathName, error) {
	// Note: We could call GlobUpspin but that might do an unnecessary Lookup.
	pat := s.AtSign(pattern)
	parsed, err := path.Parse(pat)
	if err != nil {
		return nil, err
	}
	// If it has no metacharacters, leave it alone but clean it.
	if !HasGlobChar(string(pat)) {
		return []upspin.PathName{parsed.Path()}, nil
	}
	entries, err := s.c.Glob(parsed.String())
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("no path matches %q", parsed)
	}
	names := make([]upspin.PathName, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name
	}
	return names, nil
}

// GlobAllUpspinPath processes the arguments, which should be Upspin paths,
// expanding glob patterns. It returns just the paths.
func (s *Gateway) GlobAllUpspinPath(args []string) []upspin.PathName {
	paths := make([]upspin.PathName, 0, len(args))
	for _, arg := range args {
		xs, err := s.GlobUpspinPath(arg)
		if err != nil {
			continue
		}
		paths = append(paths, xs...)
	}
	return paths
}

// expandUpspin turns the list of string arguments into Upspin path names.
// If glob is true, it "globs" and @-expands the arguments.
// Otherwise, it interprets leading @ symbols but does no other processing.
func (s *Gateway) expandUpspin(args []string, doGlob bool) []upspin.PathName {
	if doGlob {
		return s.GlobAllUpspinPath(args)
	}
	paths := make([]upspin.PathName, len(args))
	for i, arg := range args {
		paths[i] = upspin.PathName(s.AtSign(arg))
	}
	return paths
}
