import json
import sys
import numpy as np
import matplotlib.pyplot as plt
from sklearn.metrics import accuracy_score, precision_score, recall_score, f1_score
import os


def update_master_results(accuracy, master_file_path):
    if os.path.exists(master_file_path):
        with open(master_file_path, 'r') as master_file:
            master_results = json.load(master_file)
    else:
        master_results = {'accuracies': []}

    master_results['accuracies'].append(accuracy)

    with open(master_file_path, 'w') as master_file:
        json.dump(master_results, master_file, indent=4)

    return master_results['accuracies']


def create_accuracy_graph(accuracies, current_accuracy, graph_file_path):
    plt.figure()
    plt.scatter(range(len(accuracies)), accuracies, label='Past Accuracies', alpha=0.5)

    current_index = len(accuracies) - 1
    plt.scatter(current_index, current_accuracy, color='red', label='Current Analysis', zorder=5)

    plt.annotate(f'{current_accuracy:.2%}', (current_index, current_accuracy),
                 textcoords="offset points", xytext=(0,10), ha='center')

    plt.title('Accuracy Comparison Graph')
    plt.xlabel('Analysis Number')
    plt.ylabel('Accuracy')
    plt.legend()

    plt.ylim(0, 1)

    plt.savefig(graph_file_path)
    plt.close()



def plot_heatmap(data, file_name):
    y_scores = [item['model_score'] for item in data]
    
    scores_matrix = np.array(y_scores).reshape(1, -1)

    num_texts = len(data)
    fig, ax = plt.subplots(figsize=(20, 1))

    cax = ax.matshow(scores_matrix, cmap='coolwarm', interpolation='nearest')

    plt.xticks(ticks=np.linspace(0, num_texts-1, min(num_texts, 10)), labels=[f'Text {i+1}' for i in range(min(num_texts, 10))], rotation=90)
    plt.yticks([])

    plt.colorbar(cax, orientation='vertical', fraction=0.025)
    plt.title('Heat Map of Model Scores')

    plot_file_name = file_name.replace('.json', '_heatmap.png')
    plt.savefig(plot_file_name, dpi=300, bbox_inches='tight')
    plt.close()
    
    print(f"Heat map saved to: {plot_file_name}")


def analyze_data(file_name):
    print(file_name)
    input_filepath = os.path.join("..", "outputs", file_name)
    print(input_filepath)
    with open(input_filepath, 'r') as file:
        data = json.load(file)

    plot_heatmap(data, file_name) 

    y_true = np.array([item['original_score'] for item in data])
    y_scores = np.array([item['model_score'] for item in data])

    threshold = 0.5
    y_pred = np.where(y_scores >= threshold, 1, 0)

    results = {
        'accuracy': accuracy_score(y_true, y_pred),
        'precision': precision_score(y_true, y_pred, zero_division=0),
        'recall': recall_score(y_true, y_pred, zero_division=0),
        'f1_score': f1_score(y_true, y_pred, zero_division=0),
    }
    output_dir = "../outputs"
    master_file_path = os.path.join(output_dir, "master_results.json")
    all_accuracies = update_master_results(results['accuracy'], master_file_path)

    graph_file_path = os.path.join("../outputs", "accuracies_graph.png")
    create_accuracy_graph(all_accuracies, results['accuracy'], graph_file_path)
    results_file_name = input_filepath.replace('.json', '_results.json')
    results_filepath = os.path.join(output_dir, results_file_name)
    with open(results_filepath, 'w') as results_file:
        json.dump(results, results_file, indent=4)


if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python analysis.py input_file")
        sys.exit(1)

    file_name = sys.argv[1]
    analyze_data(file_name)
