import { useEffect, useState } from "react";
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  StyleSheet,
} from "react-native";

import { getTasks } from "../api/taskApi";
import { Task } from "../types/task";
import SearchBar from "../components/SearchBar";

export default function TaskScreen() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);

  const [keyword, setKeyword] = useState("");

  const fetchTasks = async () => {
    try {
      setLoading(true);

      const response = await getTasks({
        keyword,
      });

      setTasks(response.data.data);
    } catch (error) {
      console.log(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, [keyword]);

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator size="large" />
      </View>
    );
  }

  return (
    <>
      <SearchBar value={keyword} onChangeText={setKeyword} />

      <FlatList
        data={tasks}
        keyExtractor={(item) => item.id.toString()}
        renderItem={({ item }) => (
          <View style={styles.card}>
            <Text style={styles.title}>{item.title}</Text>

            <Text>{item.description}</Text>

            <Text>Status : {item.status}</Text>

            <Text>Assignee : {item.assignee}</Text>
          </View>
        )}
      />
    </>
  );
}

const styles = StyleSheet.create({
  center: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },

  card: {
    backgroundColor: "#fff",
    margin: 10,
    padding: 15,
    borderRadius: 8,
    elevation: 2,
  },

  title: {
    fontSize: 18,
    fontWeight: "bold",
    marginBottom: 5,
  },
});
