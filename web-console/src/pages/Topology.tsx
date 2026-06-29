import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Building2, Server, Thermometer, Router, ChevronRight, ChevronDown } from 'lucide-react';
import { mockDataCenters, mockEdgeAgents, mockServers, mockCoolingDevices } from '../api/mock';

interface TreeNode {
  id: string;
  name: string;
  type: 'datacenter' | 'agent' | 'server' | 'cooling' | 'network';
  status?: string;
  children?: TreeNode[];
}

function Topology() {
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set(['dc-1']));

  const { data: dataCenters } = useQuery({
    queryKey: ['dataCenters'],
    queryFn: () => mockDataCenters,
  });

  const { data: agents } = useQuery({
    queryKey: ['agents'],
    queryFn: () => mockEdgeAgents,
  });

  const { data: servers } = useQuery({
    queryKey: ['servers'],
    queryFn: () => mockServers,
  });

  const { data: coolingDevices } = useQuery({
    queryKey: ['coolingDevices'],
    queryFn: () => mockCoolingDevices,
  });

  const buildTree = (): TreeNode[] => {
    return (
      dataCenters?.map((dc) => ({
        id: dc.id,
        name: dc.name,
        type: 'datacenter' as const,
        children: agents
          ?.filter((agent) => agent.dc_id === dc.id)
          .map((agent) => ({
            id: agent.id,
            name: agent.hostname,
            type: 'agent' as const,
            status: agent.status,
            children: [
              ...(servers
                ?.filter((server) => server.agent_id === agent.id)
                .map((server) => ({
                  id: server.id,
                  name: `${server.manufacturer} ${server.model}`,
                  type: 'server' as const,
                  status: server.status,
                })) || []),
              ...(coolingDevices
                ?.filter((device) => device.agent_id === agent.id)
                .map((device) => ({
                  id: device.id,
                  name: device.name,
                  type: 'cooling' as const,
                  status: device.status,
                })) || []),
            ],
          })),
      })) || []
    );
  };

  const toggleNode = (nodeId: string) => {
    const newExpanded = new Set(expandedNodes);
    if (newExpanded.has(nodeId)) {
      newExpanded.delete(nodeId);
    } else {
      newExpanded.add(nodeId);
    }
    setExpandedNodes(newExpanded);
  };

  const getNodeIcon = (type: TreeNode['type']) => {
    switch (type) {
      case 'datacenter':
        return <Building2 className="text-blue-500" size={20} />;
      case 'agent':
        return <Router className="text-purple-500" size={20} />;
      case 'server':
        return <Server className="text-green-500" size={20} />;
      case 'cooling':
        return <Thermometer className="text-cyan-500" size={20} />;
      default:
        return <Server className="text-gray-500" size={20} />;
    }
  };

  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'online':
        return 'bg-green-500';
      case 'offline':
        return 'bg-red-500';
      case 'error':
        return 'bg-red-500';
      default:
        return 'bg-gray-400';
    }
  };

  const renderNode = (node: TreeNode, level: number = 0) => {
    const isExpanded = expandedNodes.has(node.id);
    const hasChildren = node.children && node.children.length > 0;

    return (
      <div key={node.id} className="select-none">
        <div
          className="flex items-center gap-2 py-2 px-3 rounded-lg hover:bg-gray-100 cursor-pointer"
          style={{ paddingLeft: `${level * 24 + 12}px` }}
          onClick={() => hasChildren && toggleNode(node.id)}
        >
          {hasChildren ? (
            isExpanded ? (
              <ChevronDown size={16} className="text-gray-400" />
            ) : (
              <ChevronRight size={16} className="text-gray-400" />
            )
          ) : (
            <span className="w-4" />
          )}
          {getNodeIcon(node.type)}
          <span className="flex-1 font-medium text-gray-800">{node.name}</span>
          {node.status && (
            <>
              <span className={`w-2 h-2 rounded-full ${getStatusColor(node.status)}`} />
              <span className="text-xs text-gray-500 capitalize">{node.status}</span>
            </>
          )}
        </div>
        {hasChildren && isExpanded && (
          <div>{node.children!.map((child) => renderNode(child, level + 1))}</div>
        )}
      </div>
    );
  };

  const tree = buildTree();

  return (
    <div className="space-y-6">
      <h2 className="text-2xl font-bold text-gray-800">Topology</h2>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Tree View */}
        <div className="lg:col-span-1 bg-white rounded-lg shadow">
          <div className="p-4 border-b">
            <h3 className="font-semibold text-gray-800">Infrastructure Tree</h3>
          </div>
          <div className="p-2">{tree.map((node) => renderNode(node))}</div>
        </div>

        {/* Details Panel */}
        <div className="lg:col-span-2 bg-white rounded-lg shadow">
          <div className="p-4 border-b">
            <h3 className="font-semibold text-gray-800">Topology Overview</h3>
          </div>
          <div className="p-6">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              <div className="text-center p-4 bg-blue-50 rounded-lg">
                <Building2 className="mx-auto text-blue-500 mb-2" size={32} />
                <p className="text-2xl font-bold text-gray-800">{dataCenters?.length || 0}</p>
                <p className="text-sm text-gray-500">Data Centers</p>
              </div>
              <div className="text-center p-4 bg-purple-50 rounded-lg">
                <Router className="mx-auto text-purple-500 mb-2" size={32} />
                <p className="text-2xl font-bold text-gray-800">{agents?.length || 0}</p>
                <p className="text-sm text-gray-500">Edge Agents</p>
              </div>
              <div className="text-center p-4 bg-green-50 rounded-lg">
                <Server className="mx-auto text-green-500 mb-2" size={32} />
                <p className="text-2xl font-bold text-gray-800">{servers?.length || 0}</p>
                <p className="text-sm text-gray-500">Servers</p>
              </div>
              <div className="text-center p-4 bg-cyan-50 rounded-lg">
                <Thermometer className="mx-auto text-cyan-500 mb-2" size={32} />
                <p className="text-2xl font-bold text-gray-800">{coolingDevices?.length || 0}</p>
                <p className="text-sm text-gray-500">Cooling Devices</p>
              </div>
            </div>

            <div className="mt-6 p-4 bg-gray-50 rounded-lg">
              <h4 className="font-medium text-gray-800 mb-2">Legend</h4>
              <div className="flex flex-wrap gap-4">
                <div className="flex items-center gap-2">
                  <Building2 className="text-blue-500" size={16} />
                  <span className="text-sm text-gray-600">Data Center</span>
                </div>
                <div className="flex items-center gap-2">
                  <Router className="text-purple-500" size={16} />
                  <span className="text-sm text-gray-600">Edge Agent</span>
                </div>
                <div className="flex items-center gap-2">
                  <Server className="text-green-500" size={16} />
                  <span className="text-sm text-gray-600">Server</span>
                </div>
                <div className="flex items-center gap-2">
                  <Thermometer className="text-cyan-500" size={16} />
                  <span className="text-sm text-gray-600">Cooling Device</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-green-500" />
                  <span className="text-sm text-gray-600">Online</span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-red-500" />
                  <span className="text-sm text-gray-600">Offline</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Topology;
