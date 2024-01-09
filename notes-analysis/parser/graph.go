package parser

import (
	"github.com/zoroqi/hamster/notes-analysis/store"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

func MakingGraph(docs []*store.Document) map[string]*store.GraphNode {
	nodes := map[string]*store.GraphNode{}
	for _, doc := range docs {
		nodes[doc.Id] = &store.GraphNode{
			Id:       doc.Id,
			Doc:      doc,
			OutEdges: make([]store.GraphEdge, 0),
		}
		if doc.Type == store.D_MD {
			for _, block := range doc.Blocks {
				store.WalkBlocks(block, func(bk *store.Block) bool {
					id := doc.Id + "_" + bk.Id
					nodes[id] = &store.GraphNode{
						Id:       id,
						Doc:      doc,
						OutEdges: make([]store.GraphEdge, 0),
						Block:    bk,
					}
					return true
				})
			}
		}
	}

	nameToPath := map[string][]string{}
	pathToId := map[string]string{}
	for _, v := range docs {
		nameToPath[v.Name] = append(nameToPath[v.Name], v.Path)
		pathToId[v.Path] = v.Id
	}

	for _, doc := range docs {
		if doc.Type == store.D_MD {
			for _, block := range doc.Blocks {
				store.WalkBlocks(block, func(bk *store.Block) bool {
					id := doc.Id + "_" + bk.Id
					for _, link := range bk.Links {
						if link.Type == store.L_FILE {
							filename, path := store.TargetFileName(link.Target)
							paths := nameToPath[filename]
							if len(paths) >= 1 {
								// 找到最可能的文件
							} else {
								// no file, so create doc
								newDoc := createEmptyDoc(filepath.Join("nofile", path))
								newDoc.Meta = yaml.MapSlice{
									{
										Key:   "tags",
										Value: []string{link.Alias},
									},
								}
								nodes[newDoc.Id] = &store.GraphNode{
									Id:       newDoc.Id,
									Doc:      newDoc,
									OutEdges: make([]store.GraphEdge, 0),
								}

								nameToPath[newDoc.Name] = append(nameToPath[newDoc.Name], newDoc.Path)
								pathToId[newDoc.Path] = newDoc.Id
								paths = append(paths, newDoc.Path)
							}

							nodes[id].OutEdges = append(nodes[id].OutEdges, store.GraphEdge{
								From:  id,
								To:    pathToId[paths[0]],
								Alias: link.Alias,
							})
						}
					}
					return true
				})
			}
		}
	}

	return nodes
}

func WalkGraph(nodes map[string]*store.GraphNode, walker func(from *store.GraphNode, to *store.GraphNode) error) error {
	visit := map[string]bool{}
	var walk func(node *store.GraphNode) error
	walk = func(node *store.GraphNode) error {
		if visit[node.Id] {
			return nil
		}
		visit[node.Id] = true
		walker(node, nil)
		for _, edge := range node.OutEdges {
			if err := walker(node, nodes[edge.To]); err != nil {
				return err
			}
			if err := walk(nodes[edge.To]); err != nil {
				return err
			}
		}
		return nil
	}
	for _, root := range nodes {
		if err := walk(root); err != nil {
			return err
		}
	}
	return nil
}
