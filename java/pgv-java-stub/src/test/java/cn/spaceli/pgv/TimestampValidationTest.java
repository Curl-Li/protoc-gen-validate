package cn.spaceli.pgv;

import com.google.protobuf.Duration;
import com.google.protobuf.Timestamp;
import com.google.protobuf.util.Durations;
import com.google.protobuf.util.Timestamps;
import org.junit.Test;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class TimestampValidationTest {
    @Test
    public void lessThanWorks() throws RuntimeException {
        TestException ex = new TestException(2, "time less than target");
        // Less
        ComparativeValidation.lessThan(ex, Timestamps.fromSeconds(10), Timestamps.fromSeconds(20), Timestamps.comparator());
        // Equal
        assertThatThrownBy(() -> ComparativeValidation.lessThan(ex, Timestamps.fromSeconds(10), Timestamps.fromSeconds(10), Timestamps.comparator())).isEqualTo(ex);
        // Greater
        assertThatThrownBy(() -> ComparativeValidation.lessThan(ex, Timestamps.fromSeconds(20), Timestamps.fromSeconds(10), Timestamps.comparator())).isEqualTo(ex);
    }

    @Test
    public void lessThanOrEqualsWorks() throws RuntimeException {
        TestException ex = new TestException(2, "time less than target");
        // Less
        ComparativeValidation.lessThanOrEqual(ex, Timestamps.fromSeconds(10), Timestamps.fromSeconds(20), Timestamps.comparator());
        // Equal
        ComparativeValidation.lessThanOrEqual(ex, Timestamps.fromSeconds(10), Timestamps.fromSeconds(10), Timestamps.comparator());
        // Greater
        assertThatThrownBy(() -> ComparativeValidation.lessThanOrEqual(ex, Timestamps.fromSeconds(20), Timestamps.fromSeconds(10), Timestamps.comparator())).isEqualTo(ex);
    }

    @Test
    public void greaterThanWorks() throws RuntimeException {
        TestException ex = new TestException(2, "time greater than target");
        // Less
        assertThatThrownBy(() -> ComparativeValidation.greaterThan(ex, Timestamps.fromSeconds(10), Timestamps.fromSeconds(20), Timestamps.comparator())).isEqualTo(ex);
        // Equal
        assertThatThrownBy(() -> ComparativeValidation.greaterThan(ex, Timestamps.fromSeconds(10), Timestamps.fromSeconds(10), Timestamps.comparator())).isEqualTo(ex);
        // Greater
        ComparativeValidation.greaterThan(ex, Timestamps.fromSeconds(20), Timestamps.fromSeconds(10), Timestamps.comparator());
    }

    @Test
    public void greaterThanOrEqualsWorks() throws RuntimeException {
        TestException ex = new TestException(2, "time greater than target");
        // Less
        assertThatThrownBy(() -> ComparativeValidation.greaterThanOrEqual(ex, Timestamps.fromSeconds(10), Timestamps.fromSeconds(20), Timestamps.comparator())).isEqualTo(ex);
        // Equal
        ComparativeValidation.greaterThanOrEqual(ex, Timestamps.fromSeconds(10), Timestamps.fromSeconds(10), Timestamps.comparator());
        // Greater
        ComparativeValidation.greaterThanOrEqual(ex, Timestamps.fromSeconds(20), Timestamps.fromSeconds(10), Timestamps.comparator());
    }

    @Test
    public void withinWorks() throws RuntimeException {
        TestException ex = new TestException(2, "time not in specify range");
        Timestamp when = Timestamps.fromSeconds(20);
        Duration duration = Durations.fromSeconds(5);

        // Less
        TimestampValidation.within(ex, Timestamps.fromSeconds(18), duration, when);
        TimestampValidation.within(ex, Timestamps.fromSeconds(20), duration, when);
        TimestampValidation.within(ex, Timestamps.fromSeconds(22), duration, when);

        // Equal
        TimestampValidation.within(ex, Timestamps.fromSeconds(15), duration, when);
        TimestampValidation.within(ex, Timestamps.fromSeconds(25), duration, when);

        // Greater
        assertThatThrownBy(() -> TimestampValidation.within(ex, Timestamps.fromSeconds(10), duration, when)).isEqualTo(ex);
        assertThatThrownBy(() -> TimestampValidation.within(ex, Timestamps.fromSeconds(30), duration, when)).isEqualTo(ex);
    }
}
