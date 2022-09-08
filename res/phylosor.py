from dendropy import Tree
import argparse

def get_edge_length( node ):
    """Gets edge length of node, even if not present.
    Parameters
    ----------
    node : dendropy.Node
    Returns
    -------
    float
        branch length subtending node. Will be 0 for root.
    """
    if node.edge_length is None:
        return 0
    else:
        return node.edge_length

def phylosor( tree, comA, comB ):
    """ Calculates the branch lengths of two communities, as well as the branch lengths that they share.
    Parameters
    ----------
    tree : dendropy.Tree
        phylogenetic tree containing taxa from both communities.
    comA : list
        set of taxa collected from first community.
    comB : list
        set of taxa collected from second community.
    Returns
    -------
    blA : float
        total branch length of all taxa in first community.
    blB : float
        total branch length of all taxa in second community.
    blBoth : float
        total branch length shared by both communities
    """
    blA = 0
    blB = 0
    blBoth = 0

    tree0 = tree.extract_tree()
    for i in tree0.leaf_nodes():
        if i.taxon.label in comA:
            blA += get_edge_length( i )
            for j in i.ancestor_iter():
                if getattr( j, "comA", False ):
                    break
                elif j.edge.length is not None:
                    j.comA = True
                    blA += j.edge.length
                    if getattr( j, "comB", False ):
                        blBoth += j.edge.length
        elif i.taxon.label in comB:
            blB += get_edge_length( i )
            for j in i.ancestor_iter():
                if getattr( j, "comB", False ):
                    break
                elif j.edge.length is not None:
                    j.comB = True
                    blB += j.edge.length
                    if getattr( j, "comA", False ):
                        blBoth += j.edge.length
    return blBoth / ( 0.5 * (blA + blB) )

if __name__ == "__main__":
    parser = argparse.ArgumentParser( description="Calculates phylosor metric between location pair" )
    parser.add_argument( "--tree", type=str, help="input tree", required=True )
    parser.add_argument( "--commA", nargs="+", help="taxa in community A", required=True )
    parser.add_argument( "--commB", nargs="+", help="taxa in community B", required=True )

    args = parser.parse_args()

    tree = Tree.get( path=args.tree, schema="newick" )
    result = phylosor( tree, args.commA, args.commB )
    print( result )