package merkledag

import (
	"errors"
	"strings"
)

func Hash2File(store KVStore, hash []byte, path string, hp HashPool) ([]byte, error) {
	// 从KVStore中获取节点
	nodeData, err := store.Get(hash)
	if err != nil {
		return nil, err
	}

	// 这里我们假设Node的数据是以序列化的形式存储的，我们需要反序列化它
	// 你需要根据你的实际情况来调整这部分代码
	node := deserialize(nodeData)

	// 检查节点是否为File
	if file, ok := node.(File); ok {
		return file.Bytes(), nil
	}

	// 检查节点是否为Dir
	if dir, ok := node.(Dir); ok {
		// 分割路径
		parts := strings.Split(path, "/")

		// 获取迭代器
		it := dir.It()

		// 遍历Dir中的文件/目录
		for it.Next() {
			// 如果当前项的名称与路径的第一部分匹配
			if it.Node().Name() == parts[0] {
				// 如果这是路径的最后一部分，返回项的内容
				if len(parts) == 1 {
					if file, ok := it.Node().(File); ok {
						return file.Bytes(), nil
					}
				} else {
					// 否则，递归调用Hash2File，使用路径的其余部分
					return Hash2File(store, it.Node().Hash(), strings.Join(parts[1:], "/"), hp)
				}
			}
		}
	}

	// 如果找不到路径，返回错误
	return nil, errors.New("path not found")
}

// 这是一个假设的反序列化函数，你需要根据你的实际情况来实现它
func deserialize(data []byte) Node {
	// TODO: 实现这个函数
	return nil
}
