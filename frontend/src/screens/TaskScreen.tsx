import { useEffect, useState } from "react";
import { View, FlatList, ActivityIndicator, StyleSheet } from "react-native";

import { getTasks } from "../api/taskApi";
import { Task } from "../types/task";

import SearchBar from "../components/SearchBar";
import StatusFilter from "../components/StatusFilter";
import Pagination from "../components/Pagination";
import EditTaskModal from "../components/EditTaskModal";
import TaskCard from "../components/TaskCard";

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
          <TaskCard
            task={item}
            onEdit={() => {
              setSelectedTask(item);
              setModalVisible(true);
            }}
          />
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
});
