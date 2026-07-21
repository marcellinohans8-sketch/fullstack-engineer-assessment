import { ActivityIndicator, StyleSheet, View } from "react-native";

export default function Loading() {
  return (
    <View style={styles.container} testID="loading-state">
      <ActivityIndicator size="large" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    padding: 24,
  },
});
