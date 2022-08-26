package cn.spaceli.pgv;

import com.google.re2j.Pattern;
import org.apache.commons.validator.routines.DomainValidator;
import org.apache.commons.validator.routines.EmailValidator;
import org.apache.commons.validator.routines.InetAddressValidator;

import java.net.URI;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;

/**
 * {@code StringValidation} implements PGV validation for protobuf {@code String} fields.
 */
@SuppressWarnings("WeakerAccess")
public final class StringValidation {
    private static final int UUID_DASH_1 = 8;
    private static final int UUID_DASH_2 = 13;
    private static final int UUID_DASH_3 = 18;
    private static final int UUID_DASH_4 = 23;
    private static final int UUID_LEN = 36;

    private StringValidation() {
        // Intentionally left blank.
    }

    // Defers initialization until needed and from there on we keep an object
    // reference and avoid future calls; it is safe to assume that we require
    // the instance again after initialization.
    private static class Lazy {
        static final EmailValidator EMAIL_VALIDATOR = EmailValidator.getInstance(true, true);
    }

    public static void length(final RuntimeException ex, final String value, final int expected) {
        final int actual = value.codePointCount(0, value.length());
        if (actual != expected) {
            throw ex;
        }
    }

    public static void minLength(final RuntimeException ex, final String value, final int expected) {
        final int actual = value.codePointCount(0, value.length());
        if (actual < expected) {
            throw ex;
        }
    }

    public static void maxLength(final RuntimeException ex, final String value, final int expected) {
        final int actual = value.codePointCount(0, value.length());
        if (actual > expected) {
            throw ex;
        }
    }

    public static void lenBytes(RuntimeException ex, String value, int expected) {
        if (value.getBytes(StandardCharsets.UTF_8).length != expected) {
            throw ex;
        }
    }

    public static void minBytes(RuntimeException ex, String value, int expected) {
        if (value.getBytes(StandardCharsets.UTF_8).length < expected) {
            throw ex;
        }
    }

    public static void maxBytes(RuntimeException ex, String value, int expected) {
        if (value.getBytes(StandardCharsets.UTF_8).length > expected) {
            throw ex;
        }
    }

    public static void pattern(RuntimeException ex, String value, Pattern p) {
        if (!p.matches(value)) {
            throw ex;
        }
    }

    public static void prefix(RuntimeException ex, String value, String prefix) {
        if (!value.startsWith(prefix)) {
            throw ex;
        }
    }

    public static void contains(RuntimeException ex, String value, String contains) {
        if (!value.contains(contains)) {
            throw ex;
        }
    }

    public static void notContains(RuntimeException ex, String value, String contains) {
        if (value.contains(contains)) {
            throw ex;
        }
    }


    public static void suffix(RuntimeException ex, String value, String suffix) {
        if (!value.endsWith(suffix)) {
            throw ex;
        }
    }

    public static void email(final RuntimeException ex, String value) {
        if (!value.isEmpty() && value.charAt(value.length() - 1) == '>') {
            final char[] chars = value.toCharArray();
            final StringBuilder sb = new StringBuilder();
            boolean insideQuotes = false;
            for (int i = chars.length - 2; i >= 0; i--) {
                final char c = chars[i];
                if (c == '<') {
                    if (!insideQuotes) break;
                } else if (c == '"') {
                    insideQuotes = !insideQuotes;
                }
                sb.append(c);
            }
            value = sb.reverse().toString();
        }

        if (!Lazy.EMAIL_VALIDATOR.isValid(value)) {
            throw ex;
        }
    }

    public static void address(RuntimeException ex, String value) {
        boolean validHost = isAscii(value) && DomainValidator.getInstance(true).isValid(value);
        boolean validIp = InetAddressValidator.getInstance().isValid(value);

        if (!validHost && !validIp) {
            throw ex;
        }
    }

    public static void hostName(RuntimeException ex, String value) {
        if (!isAscii(value)) {
            throw ex;
        }

        DomainValidator domainValidator = DomainValidator.getInstance(true);
        if (!domainValidator.isValid(value)) {
            throw ex;
        }
    }

    public static void ip(RuntimeException ex, String value) {
        InetAddressValidator ipValidator = InetAddressValidator.getInstance();
        if (!ipValidator.isValid(value)) {
            throw ex;
        }
    }

    public static void ipv4(RuntimeException ex, String value) {
        InetAddressValidator ipValidator = InetAddressValidator.getInstance();
        if (!ipValidator.isValidInet4Address(value)) {
            throw ex;
        }
    }

    public static void ipv6(RuntimeException ex, String value) {
        InetAddressValidator ipValidator = InetAddressValidator.getInstance();
        if (!ipValidator.isValidInet6Address(value)) {
            throw ex;
        }
    }

    public static void uri(RuntimeException ex, String value) {
        try {
            URI uri = new URI(value);
            if (!uri.isAbsolute()) {
                throw ex;
            }
        } catch (URISyntaxException e) {
            throw ex;
        }
    }

    public static void uriRef(RuntimeException ex, String value) {
        try {
            new URI(value);
        } catch (URISyntaxException e) {
            throw ex;
        }
    }

    /**
     * Validates if the given value is a UUID or GUID in RFC 4122 hyphenated
     * ({@code 00000000-0000-0000-0000-000000000000}) form; both lower and upper
     * hex digits are accepted.
     */
    public static void uuid(final RuntimeException ex, final String value) {
        final char[] chars = value.toCharArray();

        err: if (chars.length == UUID_LEN) {
            for (int i = 0; i < chars.length; i++) {
                final char c = chars[i];
                if (i == UUID_DASH_1 || i == UUID_DASH_2 || i == UUID_DASH_3 || i == UUID_DASH_4) {
                    if (c != '-') {
                        break err;
                    }
                } else if ((c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F')) {
                    break err;
                }
            }
            return;
        }

        throw ex;
    }

    private static String enquote(String value) {
        return "\"" + value + "\"";
    }

    private static boolean isAscii(final String value) {
        for (char c : value.toCharArray()) {
            if (c > 127) {
                return false;
            }
        }
        return true;
    }
}
