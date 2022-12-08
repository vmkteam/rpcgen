import Foundation

extension RPCAPI: RPCMethod {
    public var rpcMethod: String {
        switch self {
        case .catalogueFirst: return "catalogue.First"
        case .catalogueSecond: return "catalogue.Second"
        case .catalogueThird: return "catalogue.Third"
        }
    }
}

extension RPCAPI: RPCParameters {
    public var rpcParameters: [String: Any?]? {
        switch self {
        case let .catalogueFirst(groups):
            return ["groups": groups.any]

        case let .catalogueSecond(campaigns):
            return ["campaigns": campaigns.any]

        case .catalogueThird:
            return nil
        }
    }
}

public enum RPCAPI {
    /// - Returns: Bool
    case catalogueFirst(groups: [Group])

    /// - Returns: Bool
    case catalogueSecond(campaigns: [Campaign])

    /// - Returns: Campaign
    case catalogueThird

}


public struct Campaign: Codable, Hashable {
    @DecodableDefault.EmptyList
    var groups: [Group]
    @DecodableDefault.IntegerZero
    var id: Int
    init(groups: [Group], id: Int) {
        self.groups = groups
        self.id = id
    }
}

public struct Group: Codable, Hashable {
    var child: Group?
    @DecodableDefault.EmptyList
    var groups: [Group]
    @DecodableDefault.IntegerZero
    var id: Int
    @DecodableDefault.EmptyList
    var nodes: [Group]
    var sub: SubGroup
    @DecodableDefault.EmptyString
    var title: String
    init(child: Group? = nil, groups: [Group], id: Int, nodes: [Group], sub: SubGroup, title: String) {
        self.child = child
        self.groups = groups
        self.id = id
        self.nodes = nodes
        self.sub = sub
        self.title = title
    }
}

public struct SubGroup: Codable, Hashable {
    @DecodableDefault.IntegerZero
    var id: Int
    @DecodableDefault.EmptyList
    var nodes: [Group]
    @DecodableDefault.EmptyString
    var title: String
    init(id: Int, nodes: [Group], title: String) {
        self.id = id
        self.nodes = nodes
        self.title = title
    }
}

