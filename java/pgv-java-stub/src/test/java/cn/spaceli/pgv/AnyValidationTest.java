package cn.spaceli.pgv;

import com.google.protobuf.Any;
import org.junit.Test;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class AnyValidationTest {

    @Test
    public void inWorks() throws RuntimeException {
        String[] set = new String[]{"type.googleapis.com/google.protobuf.Duration"};
        TestException ex = new TestException(2, "value not in set");

        // In
        CollectiveValidation.in(ex, Any.newBuilder().setTypeUrl("type.googleapis.com/google.protobuf.Duration").build().getTypeUrl(), set);

        // Not In
        assertThatThrownBy(() -> CollectiveValidation.in(ex, Any.newBuilder().setTypeUrl("junk").build().getTypeUrl(), set)).isEqualTo(ex);
    }

    @Test
    public void notInWorks() throws RuntimeException {
        String[] set = new String[]{"type.googleapis.com/google.protobuf.Duration"};
        TestException ex = new TestException(2, "value in set");
        // In
        assertThatThrownBy(() -> CollectiveValidation.notIn(ex, Any.newBuilder().setTypeUrl("type.googleapis.com/google.protobuf.Duration").build().getTypeUrl(), set)).isEqualTo(ex);

        // Not In
        CollectiveValidation.notIn(ex, Any.newBuilder().setTypeUrl("junk").build().getTypeUrl(), set);
    }
}
