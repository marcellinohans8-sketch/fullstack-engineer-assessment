import { useEffect, useState } from "react";
import {
  View,
  Text,
  FlatList,
  ActivityIndicator,
  StyleSheet,
  Button,
} from "react-native";

import { getTasks } from "../api/taskApi";
import { Task } from "../types/task";

import SearchBar from "../components/SearchBar";
import StatusFilter from "../components/StatusFilter";
import Pagination from "../components/Pagination";
import EditTaskModal from "../components/EditTaskModal";

export default function TaskScreen() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);

  const [keyword, setKeyword] = useState("");
  const [status, setStatus] = useState("");

  const [page, setPage] = useState(1);
  const [limit] = useState(5);
  const [total, setTotal] = useState(0);

  const [selectedTask, setSelectedTask] = useState<Task | null>(null);
  const [modalVisible, setModalVisible] = useState(false);

  const fetchTasks = async () => {
    try {
      setLoading(true);

      const response = await getTasks({
        keyword,
        status,
        page,
        limit,
      });

      setTasks(response.data.data);
      setTotal(response.data.pagination.total);
    } catch (error) {
      console.log(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, [keyword, status, page]);

  useEffect(() => {
    setPage(1);
  }, [keyword, status]);

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

      <StatusFilter value={status} onChange={setStatus} />

      <Pagination
        page={page}
        total={total}
        limit={limit}
        onPrevious={() => setPage((prev) => Math.max(prev - 1, 1))}
        onNext={() => {
          const totalPages = Math.ceil(total / limit);

          if (page < totalPages) {
            setPage((prev) => prev + 1);
          }
        }}
      />

      <FlatList
        data={tasks}
        keyExtractor={(item) => item.id.toString()}
        renderItem={({ item }) => (
          <View style={styles.card}>
            <Text style={styles.title}>{item.title}</Text>

            <Text>{item.description}</Text>

            <Text>Status: {item.status}</Text>

            <Text>Assignee: {item.assignee}</Text>

            <View style={styles.buttonContainer}>
              <Button
                title="Edit"
                onPress={() => {
                  setSelectedTask(item);
                  setModalVisible(true);
                }}
              />
            </View>
          </View>
        )}
      />

      <EditTaskModal
        visible={modalVisible}
        task={selectedTask}
        onClose={() => setModalVisible(false)}
        onSuccess={() => {
          setModalVisible(false);
          fetchTasks();
        }}
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
    marginHorizontal: 10,
    marginVertical: 6,
    padding: 15,
    borderRadius: 8,
    elevation: 2,
  },

  title: {
    fontSize: 18,
    fontWeight: "bold",
    marginBottom: 5,
  },

  buttonContainer: {
    marginTop: 10,
  },
});
