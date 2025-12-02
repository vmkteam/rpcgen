/// Code generated from jsonrpc schema by rpcgen v2.7.0; DO NOT EDIT.

import Foundation

extension RPCAPI: RPCMethod {
    public var rpcMethod: String {
        switch self {
        case .batch(let requests): return requests.compactMap { $0.rpcMethod }.joined(separator: ",")
        case .catalogueFirst: return "catalogue.First"
        case .catalogueSecond: return "catalogue.Second"
        case .catalogueThird: return "catalogue.Third"
        }
    }
}

extension RPCAPI: RPCParameters {
    public var rpcParameters: [String: Any?]? {
        switch self {
        case .batch:
              return nil
        case let .catalogueFirst(groups):
            return ["groups": groups.any]

        case let .catalogueSecond(campaigns):
            return ["campaigns": campaigns.any]

        case .catalogueThird:
            return nil
        }
    }
}

public enum RPCAPI: Codable, Hashable {
    /// Make batch requests.
    case batch(requests: [RPCAPI])

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

