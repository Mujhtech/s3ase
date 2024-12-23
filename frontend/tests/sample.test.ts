import {
  describe,
  test,
  expect,
  beforeAll,
  afterAll,
  beforeEach,
} from "vitest";

// Sample functions to test
const add = (a: number, b: number) => a + b;
const fetchUser = async (id: number) => {
  return { id, name: "John Doe" };
};

describe("Sample Test Suite", () => {
  beforeAll(() => {
    console.log("Running before all tests");
  });

  afterAll(() => {
    console.log("Running after all tests");
  });

  beforeEach(() => {
    console.log("Running before each test");
  });

  test("basic equality test", () => {
    expect(2 + 2).toBe(4);
  });

  test("testing add function", () => {
    expect(add(2, 3)).toBe(5);
    expect(add(-1, 1)).toBe(0);
  });

  test("testing objects", () => {
    const obj = { name: "John", age: 30 };
    expect(obj).toEqual({ name: "John", age: 30 });
    expect(obj).toHaveProperty("name");
  });

  test("testing arrays", () => {
    const arr = [1, 2, 3];
    expect(arr).toContain(2);
    expect(arr).toHaveLength(3);
  });

  test("async test", async () => {
    const user = await fetchUser(1);
    expect(user).toEqual({ id: 1, name: "John Doe" });
  });

  describe("Nested group of tests", () => {
    test("truthiness", () => {
      expect(true).toBeTruthy();
      expect(false).toBeFalsy();
      expect(null).toBeNull();
    });

    test("numbers", () => {
      expect(2).toBeGreaterThan(1);
      expect(3).toBeLessThan(5);
    });
  });
});
