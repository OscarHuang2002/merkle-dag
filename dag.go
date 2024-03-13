package merkledag

import (
	"bytes"
	"fmt"
	"hash"
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

// Add 函数将 Node 中的数据保存在 KVStore 中，并计算出 Merkle Root
func Add(store KVStore, node Node, h hash.Hash) ([]byte, error) {
	var data []byte
	var err error

	// 根据 Node 的类型进行不同的操作
	switch n := node.(type) {
	case File: // 如果 Node 是文件
		data = n.Bytes() // 获取文件的内容
	case Dir: // 如果 Node 是目录
		var hashes [][]byte
		it := n.It()    // 获取目录的迭代器
		for it.Next() { // 遍历目录中的每个节点
			childHash, err := Add(store, it.Node(), h) // 递归调用 Add 函数
			if err != nil {
				return nil, err
			}
			hashes = append(hashes, childHash) // 将子节点的哈希值添加到哈希值列表中
		}
		data = bytes.Join(hashes, []byte{}) // 将所有子节点的哈希值连接起来
	default:
		return nil, fmt.Errorf("unknown node type") // 如果 Node 的类型未知，返回错误
	}

	h.Reset()     // 重置哈希函数
	h.Write(data) // 计算数据的哈希值
	hash := h.Sum(nil)

	err = store.Put(hash, data) // 将哈希值和数据存储在 KVStore 中
	if err != nil {
		return nil, err
	}

	return hash, nil // 返回 Merkle Root
}
