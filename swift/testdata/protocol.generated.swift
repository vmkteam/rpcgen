/// Code generated from jsonrpc schema by rpcgen v2.4.5; DO NOT EDIT.

import Foundation

protocol CatalogueNetworking {
    func catalogueFirst(groups: [Group]) async -> Result<Bool, RpcError>
    func catalogueSecond(campaigns: [Campaign]) async -> Result<Bool, RpcError>
    func catalogueThird() async -> Result<Campaign, RpcError>
}

extension Networking: CatalogueNetworking {
    func catalogueFirst(groups: [Group]) async -> Result<Bool, RpcError> {
        await request(.catalogueFirst(groups: groups))
    }

    func catalogueSecond(campaigns: [Campaign]) async -> Result<Bool, RpcError> {
        await request(.catalogueSecond(campaigns: campaigns))
    }

    func catalogueThird() async -> Result<Campaign, RpcError> {
        await request(.catalogueThird())
    }
}
