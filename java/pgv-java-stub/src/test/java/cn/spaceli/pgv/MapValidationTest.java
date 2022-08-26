package cn.spaceli.pgv;

import org.junit.Test;

import java.util.HashMap;
import java.util.Map;

import static org.assertj.core.api.Assertions.assertThatThrownBy;

public class MapValidationTest {
    @Test
    public void minWorks() throws RuntimeException {
        TestException ex = new TestException(2, "key-values not enough");
        Map<String,String> map = new HashMap<>();
        map.put("1", "ONE");
        map.put("2", "TWO");

        // Equals
        MapValidation.min(ex, map, 2);
        // Not Equals
        assertThatThrownBy(() -> MapValidation.min(ex, map, 3)).isEqualTo(ex);
    }

    @Test
    public void maxWorks() throws RuntimeException {
        TestException ex = new TestException(2, "key-values too much");
        Map<String,String> map = new HashMap<>();
        map.put("1", "ONE");
        map.put("2", "TWO");

        // Equals
        MapValidation.max(ex, map, 2);
        // Not Equals
        assertThatThrownBy(() -> MapValidation.max(ex, map, 1)).isEqualTo(ex);
    }

    @Test
    public void noSparseWorks() throws RuntimeException {
        TestException ex = new TestException(2, "null value exists");
        Map<String,String> map = new HashMap<>();
        map.put("1", "ONE");
        map.put("2", null);

        // Sparse Map
        assertThatThrownBy(() -> MapValidation.noSparse(ex, map)).isInstanceOf(RuntimeException.class);
    }
}
