<?php
/**
* PHP RPC Client by rpcgen
*/

namespace JsonRpcClient;

	use EazyJsonRpc\BaseJsonRpcClient;
    use EazyJsonRpc\BaseJsonRpcException;
    use GuzzleHttp\Exception\GuzzleException;
    use JsonMapper_Exception;

	
	/** Campaign **/
	class Campaign { 
	    /**
	    * @var Group[]
	    */
	    public $group;
	    /**
	    * @var int
	    */
	    public $id;
	}
	
	/** Group **/
	class Group { 
	    /**
	    * @var Group|null
	    */
	    public $child;
	    /**
	    * @var Group[]
	    */
	    public $group;
	    /**
	    * @var int
	    */
	    public $id;
	    /**
	    * @var Group[]
	    */
	    public $nodes;
	    /**
	    * @var SubGroup
	    */
	    public $sub;
	    /**
	    * @var string
	    */
	    public $title;
	}
	
	/** SubGroup **/
	class SubGroup { 
	    /**
	    * @var int
	    */
	    public $id;
	    /**
	    * @var string
	    */
	    public $title;
	}
	


    /**
     * RpcClient
     */
 	class RpcClient extends BaseJsonRpcClient {
		
		/**
		* <catalogue.First> RPC method
		* 
		* @param Group[] $groups
		* @param bool $isNotification [optional] set to true if call is notification
		* @return bool
		* @throws BaseJsonRpcException
		* @throws GuzzleException
		* @throws JsonMapper_Exception
		**/
		public function catalogue_First( array $groups, $isNotification = false ): bool {
			return $this->call( __FUNCTION__, 'bool', [ 'groups' => $groups, ], $this->getRequestId( $isNotification ) );
		}
		
		/**
		* <catalogue.Second> RPC method
		* 
		* @param Campaign[] $campaigns
		* @param bool $isNotification [optional] set to true if call is notification
		* @return bool
		* @throws BaseJsonRpcException
		* @throws GuzzleException
		* @throws JsonMapper_Exception
		**/
		public function catalogue_Second( array $campaigns, $isNotification = false ): bool {
			return $this->call( __FUNCTION__, 'bool', [ 'campaigns' => $campaigns, ], $this->getRequestId( $isNotification ) );
		}
		
		/**
		* <catalogue.Third> RPC method
		* 
		* @param bool $isNotification [optional] set to true if call is notification
		* @return object
		* @throws BaseJsonRpcException
		* @throws GuzzleException
		* @throws JsonMapper_Exception
		**/
		public function catalogue_Third( $isNotification = false ): object {
			return $this->call( __FUNCTION__, 'object', [ ], $this->getRequestId( $isNotification ) );
		}
		

        /**
         * Get Instance
         * @param $url string RPC server url
         * @return RpcClient
         */
        public static function GetInstance( $url ): RpcClient {
            return new self( $url );
        }
 	}
