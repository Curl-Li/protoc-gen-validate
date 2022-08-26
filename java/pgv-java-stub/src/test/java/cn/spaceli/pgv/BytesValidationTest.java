package cn.spaceli.pgv;

import com.google.protobuf.ByteString;
import com.google.re2j.Pattern;
import org.junit.Test;

import java.net.InetAddress;
import java.net.UnknownHostException;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class BytesValidationTest {
    @Test
    public void lengthWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Short
        assertThatThrownBy(() -> BytesValidation.length(ex, ByteString.copyFromUtf8("ñįö"), 8)).isEqualTo(ex);
        // Same
        BytesValidation.length(ex, ByteString.copyFromUtf8("ñįöxx"), 8);
        // Long
        assertThatThrownBy(() -> BytesValidation.length(ex, ByteString.copyFromUtf8("ñįöxxxx"), 8)).isEqualTo(ex);
    }

    @Test
    public void minLengthWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Short
        assertThatThrownBy(() -> BytesValidation.minLength(ex, ByteString.copyFromUtf8("ñįö"), 8)).isEqualTo(ex);
        // Same
        BytesValidation.minLength(ex, ByteString.copyFromUtf8("ñįöxx"), 8);
        // Long
        BytesValidation.minLength(ex, ByteString.copyFromUtf8("ñįöxxxx"), 8);
    }

    @Test
    public void maxLengthWorks() throws RuntimeException {
        TestException ex = new TestException(1, "length not valid");
        // Short
        BytesValidation.maxLength(ex, ByteString.copyFromUtf8("ñįö"), 8);
        // Same
        BytesValidation.maxLength(ex, ByteString.copyFromUtf8("ñįöxx"), 8);
        // Long
        assertThatThrownBy(() -> BytesValidation.maxLength(ex, ByteString.copyFromUtf8("ñįöxxxx"), 8)).isEqualTo(ex);
    }

    @Test
    public void patternWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't valid");
        Pattern p = Pattern.compile("^[\\x00-\\x7F]+$");
        // Match
        BytesValidation.pattern(ex, ByteString.copyFromUtf8("aaabbb"), p); // non-empty, ASCII byte sequence
        // No Match
        assertThatThrownBy(() -> BytesValidation.pattern(ex, ByteString.copyFromUtf8("aaañbbb"), p)).isEqualTo(ex);
    }

    @Test
    public void prefixWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't valid");
        // Match
        BytesValidation.prefix(ex, ByteString.copyFromUtf8("Hello World"), "Hello".getBytes());
        // No Match
        assertThatThrownBy(() -> BytesValidation.prefix(ex, ByteString.copyFromUtf8("Hello World"), "Bananas".getBytes())).isEqualTo(ex);
    }

    @Test
    public void containsWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string isn't contains target substring");
        // Match
        BytesValidation.contains(ex, ByteString.copyFromUtf8("Hello World"), "o W".getBytes());
        // No Match
        assertThatThrownBy(() -> BytesValidation.contains(ex, ByteString.copyFromUtf8("Hello World"), "Bananas".getBytes())).isEqualTo(ex);
    }

    @Test
    public void suffixWorks() throws RuntimeException {
        TestException ex = new TestException(1, "string doesn't have target suffix");
        // Match
        BytesValidation.suffix(ex, ByteString.copyFromUtf8("Hello World"), "World".getBytes());
        // No Match
        assertThatThrownBy(() -> BytesValidation.suffix(ex, ByteString.copyFromUtf8("Hello World"), "Bananas".getBytes())).isEqualTo(ex);
    }

    @Test
    public void ipWorks() throws RuntimeException, UnknownHostException {
        TestException ex = new TestException(1, "string isn't ip address");
        // Match
        BytesValidation.ip(ex, ByteString.copyFrom(InetAddress.getByName("192.168.0.1").getAddress()));
        BytesValidation.ip(ex, ByteString.copyFrom(InetAddress.getByName("fe80::3").getAddress()));
        // No Match
        assertThatThrownBy(() -> BytesValidation.ip(ex, ByteString.copyFromUtf8("BANANAS!"))).isEqualTo(ex);
    }

    @Test
    public void ipV4Works() throws RuntimeException, UnknownHostException {
        TestException ex = new TestException(1, "string isn't ipv4 address");
        // Match
        BytesValidation.ipv4(ex, ByteString.copyFrom(InetAddress.getByName("192.168.0.1").getAddress()));
        // No Match
        assertThatThrownBy(() -> BytesValidation.ipv4(ex, ByteString.copyFrom(InetAddress.getByName("fe80::3").getAddress()))).isEqualTo(ex);
        assertThatThrownBy(() -> BytesValidation.ipv4(ex, ByteString.copyFromUtf8("BANANAS!"))).isEqualTo(ex);
    }

    @Test
    public void ipV6Works() throws RuntimeException, UnknownHostException {
        TestException ex = new TestException(1, "string isn't ipv6 address");
        // Match
        BytesValidation.ipv6(ex, ByteString.copyFrom(InetAddress.getByName("fe80::3").getAddress()));
        // No Match
        assertThatThrownBy(() -> BytesValidation.ipv6(ex, ByteString.copyFrom(InetAddress.getByName("192.168.0.1").getAddress()))).isEqualTo(ex);
        assertThatThrownBy(() -> BytesValidation.ipv6(ex, ByteString.copyFromUtf8("BANANAS!"))).isEqualTo(ex);
    }

    @Test
    public void inWorks() throws RuntimeException {
        TestException ex = new TestException(2, "value not in set");
        ByteString[] set = new ByteString[]{ByteString.copyFromUtf8("foo"), ByteString.copyFromUtf8("bar")};
        // In
        CollectiveValidation.in(ex, ByteString.copyFromUtf8("foo"), set);
        // Not In
        assertThatThrownBy(() -> CollectiveValidation.in(ex, ByteString.copyFromUtf8("baz"), set)).isEqualTo(ex);
    }

    @Test
    public void notInWorks() throws RuntimeException {
        TestException ex = new TestException(1, "value in set");
        ByteString[] set = new ByteString[]{ByteString.copyFromUtf8("foo"), ByteString.copyFromUtf8("bar")};
        // In
        assertThatThrownBy(() -> CollectiveValidation.notIn(ex, ByteString.copyFromUtf8("foo"), set)).isEqualTo(ex);
        // Not In
        CollectiveValidation.notIn(ex, ByteString.copyFromUtf8("baz"), set);
    }
}
