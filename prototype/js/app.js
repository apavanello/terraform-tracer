// Mock de dados para simular o backend analisando arquivos Terraform
const mockNodes = [
    { id: 1, label: 'module.vpc', group: 'module', title: 'Módulo VPC base' },
    { id: 2, label: 'aws_vpc.main', group: 'resource' },
    { id: 3, label: 'aws_subnet.public', group: 'resource' },
    { id: 4, label: 'module.rds', group: 'module' },
    { id: 5, label: 'aws_db_instance.db', group: 'resource' },
    { id: 6, label: 'aws_security_group.db_sg', group: 'resource' },
    { id: 7, label: 'module.eks', group: 'module' },
    { id: 8, label: 'aws_eks_cluster.main', group: 'resource' }
];

const mockEdges = [
    { from: 1, to: 2, arrows: 'to' },
    { from: 1, to: 3, arrows: 'to' },
    { from: 4, to: 5, arrows: 'to' },
    { from: 4, to: 6, arrows: 'to' },
    { from: 7, to: 8, arrows: 'to' },
    // dependências cross-module
    { from: 5, to: 3, arrows: 'to', dashes: true, label: 'depends_on' },
    { from: 6, to: 2, arrows: 'to', dashes: true, label: 'vpc_id' },
    { from: 8, to: 3, arrows: 'to', dashes: true, label: 'subnet_ids' }
];

const mockCode = {
    5: `resource "aws_db_instance" "db" {
  identifier        = var.db_name
  engine            = "postgres"
  instance_class    = var.instance_class
  allocated_storage = var.allocated_storage

  vpc_security_group_ids = [aws_security_group.db_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
}`
};

// UI Elements
const sidebarProps = document.getElementById('props-sidebar');
const btnCloseProps = document.getElementById('btn-close-props');
const btnLoad = document.querySelector('.btn-load');
const loadingOverlay = document.getElementById('loading-overlay');
const btnFit = document.getElementById('btn-fit');

// Inicialização da Rede
function initNetwork() {
    const container = document.getElementById('network-container');
    
    // Configurações do Vis.js para estética moderna
    const options = {
        nodes: {
            shape: 'box',
            margin: 12,
            borderWidth: 2,
            borderWidthSelected: 3,
            color: {
                border: '#e5e7eb',
                background: '#ffffff',
                highlight: {
                    border: '#a30eff',
                    background: '#f3e8ff'
                },
                hover: {
                    border: '#d8b4fe',
                    background: '#faf5ff'
                }
            },
            font: {
                color: '#1f2937',
                face: 'Inter',
                size: 14,
                multi: true,
                bold: { size: 14, color: '#1f2937' }
            },
            shadow: {
                enabled: true,
                color: 'rgba(0,0,0,0.05)',
                size: 10,
                x: 0, y: 4
            }
        },
        edges: {
            width: 2,
            color: {
                color: '#cbd5e1',
                highlight: '#a30eff',
                hover: '#94a3b8'
            },
            smooth: {
                type: 'cubicBezier',
                forceDirection: 'horizontal',
                roundness: 0.6
            },
            font: {
                size: 10,
                color: '#64748b',
                face: 'Inter',
                align: 'middle'
            }
        },
        groups: {
            module: {
                color: { border: '#3b82f6', background: '#eff6ff' },
                font: { color: '#1d4ed8' },
                shape: 'database' // Apenas um shape diferente
            },
            resource: {
                shape: 'box',
                color: { border: '#e2e8f0', background: '#ffffff' }
            }
        },
        layout: {
            hierarchical: {
                enabled: true,
                direction: 'LR',
                sortMethod: 'directed',
                nodeSpacing: 150,
                levelSeparation: 250
            }
        },
        physics: false,
        interaction: {
            hover: true,
            tooltipDelay: 200,
            zoomView: true
        }
    };

    const data = {
        nodes: new vis.DataSet(mockNodes),
        edges: new vis.DataSet(mockEdges)
    };

    const network = new vis.Network(container, data, options);

    // Eventos 
    network.on('selectNode', function (params) {
        const nodeId = params.nodes[0];
        const nodeData = data.nodes.get(nodeId);
        openPropertiesPane(nodeData);
    });

    network.on('deselectNode', function () {
        sidebarProps.classList.add('closed');
    });

    // Fit Graph
    btnFit.addEventListener('click', () => {
        network.fit({ animation: { duration: 500, easingFunction: 'easeInOutQuad' } });
    });
}

function openPropertiesPane(nodeData) {
    sidebarProps.classList.remove('closed');
    document.getElementById('prop-title').innerText = nodeData.label;
    document.getElementById('prop-type').innerText = nodeData.group === 'module' ? 'Terraform Module' : 'Terraform Resource';
    document.getElementById('prop-module').innerText = nodeData.label.split('.')[0] || 'root';
    
    // Atualiza código mock
    const codeElem = document.getElementById('prop-code');
    codeElem.innerText = mockCode[nodeData.id] || `# Código não disponível para ${nodeData.label}\n\n# Clique no nó aws_db_instance.db para ver um exemplo.`;
}

// Interações da UI
btnCloseProps.addEventListener('click', () => {
    sidebarProps.classList.add('closed');
});

// Abas de Detalhes
document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', (e) => {
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
        
        const target = e.target.getAttribute('data-target');
        e.target.classList.add('active');
        document.getElementById(target).classList.add('active');
    });
});

// Simulação de Loading 
btnLoad.addEventListener('click', () => {
    loadingOverlay.style.display = 'flex';
    // forçar reflow
    loadingOverlay.offsetWidth;
    loadingOverlay.style.opacity = '1';
    
    setTimeout(() => {
        loadingOverlay.style.opacity = '0';
        setTimeout(() => {
            loadingOverlay.style.display = 'none';
        }, 300);
    }, 1500);
});

// Tree toggle visual
document.querySelectorAll('.tree-item-label').forEach(item => {
    item.addEventListener('click', (e) => {
        const parent = e.currentTarget.parentElement;
        if(parent.classList.contains('folder')) {
            parent.classList.toggle('open');
        }
    });
});

// Inicia
initNetwork();
