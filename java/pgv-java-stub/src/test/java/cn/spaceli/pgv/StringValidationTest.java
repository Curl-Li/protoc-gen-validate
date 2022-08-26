package cn.spaceli.pgv;

import com.google.re2j.Pattern;
import org.junit.Test;
import static cn.spaceli.pgv.StringValidation.uuid;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class StringValidationTest {
    private static String repeat(final char c, final int n) {
        return new String(new char[n]).replace('\0', c);
    }

    @Test
    public void inWorks() throws RuntimeException {
        String[] set = new String[]{"foo", "bar", "🙈"};
        TestException ex = new TestException(2, "value not in set");
        // In
        CollectiveValidation.in(ex, "🙈", set);
        // Not In
        assertThatThrownBy(() -> CollectiveValidation.in(ex, "baz", set)).isEqualTo(ex);
    }

    @Test
    public void notInWorks() throws RuntimeException {
        String[] set = new String[]{"foo", "bar", "🙈"};
        // In
        TestException ex = new TestException(1, "value in set");
        assertThatThrownBy(() -> CollectiveValidation.notIn(ex, "🙈", set)).isEqualTo(ex);
        // Not In
        CollectiveValidation.notIn(ex, "baz", set);
    }

    @Test
    public void lengthWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Short
        assertThatThrownBy(() -> StringValidation.length(ex, "🙈", 2)).isEqualTo(ex);
        // Same
        StringValidation.length(ex, "🙈🙉", 2);
        // Long
        assertThatThrownBy(() -> StringValidation.length(ex, "🙈🙉🙊", 2)).isEqualTo(ex);
    }

    @Test
    public void minLengthWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Short
        assertThatThrownBy(() -> StringValidation.minLength(ex, "🙈", 2)).isEqualTo(ex);
        // Same
        StringValidation.minLength(ex, "🙈🙉", 2);
        // Long
        StringValidation.minLength(ex, "🙈🙉🙊", 2);
    }

    @Test
    public void maxLengthWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Short
        StringValidation.maxLength(ex, "🙈", 2);
        // Same
        StringValidation.maxLength(ex, "🙈🙉", 2);
        // Long
        assertThatThrownBy(() -> StringValidation.maxLength(ex, "🙈🙉🙊", 2)).isEqualTo(ex);
    }

    @Test
    public void lengthBytesWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string bytes length not valid");
        // Short
        assertThatThrownBy(() -> StringValidation.lenBytes(ex, "ñįö", 8)).isEqualTo(ex);
        assertThatThrownBy(() -> StringValidation.lenBytes(ex, "ñįö", 8)).isEqualTo(ex);
        // Same
        StringValidation.lenBytes(ex, "ñįöxx", 8);
        // Long
        assertThatThrownBy(() -> StringValidation.lenBytes(ex, "ñįöxxxx", 8)).isEqualTo(ex);
    }

    @Test
    public void minBytesWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string bytes length not valid");
        // Short
        assertThatThrownBy(() -> StringValidation.minBytes(ex, "ñįö", 8)).isEqualTo(ex);
        // Same
        StringValidation.minBytes(ex, "ñįöxx", 8);
        StringValidation.minBytes(ex, "你好", 4);
        // Long
        StringValidation.minBytes(ex, "ñįöxxxx", 8);
    }

    @Test
    public void maxBytesWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string bytes length not valid");
        // Short
        StringValidation.maxBytes(ex, "ñįö", 8);
        // Same
        StringValidation.maxBytes(ex, "ñįöxx", 8);
        // Long
        assertThatThrownBy(() -> StringValidation.maxBytes(ex, "ñįöxxxx", 8)).isEqualTo(ex);
    }

    @Test
    public void patternWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't valid");
        Pattern p = Pattern.compile("a*b*");
        // Match
        StringValidation.pattern(ex, "aaabbb", p);
        // No Match
        assertThatThrownBy(() -> StringValidation.pattern(ex, "aaabbbccc", p)).isEqualTo(ex);
    }

    @Test
    public void patternWorks2() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't valid");
        Pattern p = Pattern.compile("\\* \\\\ \\w");
        // Match
        StringValidation.pattern(ex, "* \\ x", p);
    }

    @Test
    public void prefixWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't valid");
        // Match
        StringValidation.prefix(ex, "Hello World", "Hello");
        // No Match
        assertThatThrownBy(() -> StringValidation.prefix(ex, "Hello World", "Bananas")).isEqualTo(ex);
    }

    @Test
    public void containsWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't contains target substring");
        // Match
        StringValidation.contains(ex, "Hello World", "o W");
        // No Match
        assertThatThrownBy(() -> StringValidation.contains(ex, "Hello World", "Bananas")).isEqualTo(ex);
    }

    @Test
    public void notContainsWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string contains target substring");
        // Match
        StringValidation.notContains(ex, "Hello World", "Bananas");
        // No Match
        assertThatThrownBy(() -> StringValidation.notContains(ex, "Hello World", "o W")).isEqualTo(ex);
    }

    @Test
    public void suffixWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string doesn't have target suffix");
        // Match
        StringValidation.suffix(ex, "Hello World", "World");
        // No Match
        assertThatThrownBy(() -> StringValidation.suffix(ex, "Hello World", "Bananas")).isEqualTo(ex);
    }

    @Test
    public void emailWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't email address");
        // Match
        StringValidation.email(ex, "foo@bar.com");
        StringValidation.email(ex, "John Smith <foo@bar.com>");
        StringValidation.email(ex, "John Doe <john.\"we<i<>r>do\".doe@example.com>");
        StringValidation.email(ex, "john@foo.africa");
        // No Match
        assertThatThrownBy(() -> StringValidation.email(ex, "bar.bar.bar")).isEqualTo(ex);
        assertThatThrownBy(() -> StringValidation.email(ex, "John Doe <john.doe@example.com")).isEqualTo(ex);
        assertThatThrownBy(() -> StringValidation.email(ex, "John Doe <john.doe@example.com> ")).isEqualTo(ex);
    }

    @Test
    public void hostNameWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't hostname");
        // Match
        StringValidation.hostName(ex, "google.com");
        // No Match
        assertThatThrownBy(() -> StringValidation.hostName(ex, "bananas.bananas")).isEqualTo(ex);
        assertThatThrownBy(() -> StringValidation.hostName(ex, "你好.com")).isEqualTo(ex);
    }

    @Test
    public void addressWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't network address");
        // Match Hostname
        StringValidation.address(ex, "google.com");
        StringValidation.address(ex, "images.google.com");
        // Match IP
        StringValidation.address(ex, "127.0.0.1");
        StringValidation.address(ex, "fe80::3");

        // No Match
        assertThatThrownBy(() -> StringValidation.address(ex, "bananas.bananas")).isEqualTo(ex);
        assertThatThrownBy(() -> StringValidation.address(ex, "你好.com")).isEqualTo(ex);
        assertThatThrownBy(() -> StringValidation.address(ex, "ff::fff::0b")).isEqualTo(ex);
    }

    @Test
    public void ipWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't ip address");
        // Match
        StringValidation.ip(ex, "192.168.0.1");
        StringValidation.ip(ex, "fe80::3");
        // No Match
        assertThatThrownBy(() -> StringValidation.ip(ex, "999.999.999.999")).isEqualTo(ex);
    }

    @Test
    public void ipV4Works() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't ipv4 address");
        // Match
        StringValidation.ipv4(ex, "192.168.0.1");
        // No Match
        assertThatThrownBy(() -> StringValidation.ipv4(ex, "fe80::3")).isEqualTo(ex);
        assertThatThrownBy(() -> StringValidation.ipv4(ex, "999.999.999.999")).isEqualTo(ex);
    }

    @Test
    public void ipV6Works() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't ipv6 address");
        // Match
        StringValidation.ipv6(ex, "fe80::3");
        // No Match
        assertThatThrownBy(() -> StringValidation.ipv6(ex, "192.168.0.1")).isEqualTo(ex);
        assertThatThrownBy(() -> StringValidation.ipv6(ex, "999.999.999.999")).isEqualTo(ex);
    }

    @Test
    public void uriWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't uri address");
        // Match
        StringValidation.uri(ex, "ftp://ftp.is.co.za/rfc/rfc1808.txt");
        StringValidation.uri(ex, "http://www.ietf.org/rfc/rfc2396.txt");
        StringValidation.uri(ex, "ldap://[2001:db8::7]/c=GB?objectClass?one");
        StringValidation.uri(ex, "mailto:John.Doe@example.com");
        StringValidation.uri(ex, "news:comp.infosystems.www.servers.unix");
        StringValidation.uri(ex, "telnet://192.0.2.16:80/");
        StringValidation.uri(ex, "urn:oasis:names:specification:docbook:dtd:xml:4.1.2");
        StringValidation.uri(ex, "tel:+1-816-555-1212");
        // No Match
        assertThatThrownBy(() -> StringValidation.uri(ex, "server/resource")).isEqualTo(ex);
        assertThatThrownBy(() -> StringValidation.uri(ex, "this is not a uri")).isEqualTo(ex);
    }

    @Test
    public void uriRefWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't uri path");
        // Match
        StringValidation.uriRef(ex, "server/resource");
        // No Match
        assertThatThrownBy(() -> StringValidation.uri(ex, "this is not a uri")).isEqualTo(ex);
    }

    @Test
    public void uuidWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't a valid uuid");
        // We use this to generate UUIDs for all valid hex digits, so:
        // 00000000-0000…, 11111111-1111…, …, FFFFFFFF-FFFF…
        char[] chars = "0123456789abcdefABCDEF".toCharArray();

        // Match
        for (char c : chars) {
            final String s4 = repeat(c, 4);
            uuid(ex, repeat(c, 8) + '-' + s4 + '-' + s4 + '-' + s4 + '-' + repeat(c, 12));
        }

        // No Match
        assertThatThrownBy(() -> uuid(ex, "00000000-0000-0000-0000-00000000000g")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "00000000-0000_0000-0000-000000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "00000000-000000000-0000-00000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "00000000-000000000-0000-0000000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "0000000-00000-0000-0000-000000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "000000000-000-0000-0000-000000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "00000000-000-00000-0000-000000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "00000000-00000-000-0000-000000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "00000000-0000-000-00000-000000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "00000000-0000-00000-000-000000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "00000000-0000-0000-000-0000000000000")).isEqualTo(ex);
        assertThatThrownBy(() -> uuid(ex, "00000000-0000-0000-00000-00000000000")).isEqualTo(ex);
    }
}
