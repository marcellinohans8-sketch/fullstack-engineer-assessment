import { fireEvent, render } from "@testing-library/react-native";
import { describe, expect, it, jest } from "@jest/globals";

import SearchBar from "../SearchBar";

describe("SearchBar", () => {
  it("calls onChangeText when user types a keyword", async () => {
    const onChangeText = jest.fn();
    const { getByTestId } = await render(
      <SearchBar value="" onChangeText={onChangeText} />
    );

    fireEvent.changeText(getByTestId("task-search-input"), "billing");

    expect(onChangeText).toHaveBeenCalledWith("billing");
  });
});
